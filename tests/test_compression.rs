mod common;

use common::*;
use std::fs;
use tempfile::TempDir;

/// Test compression with -C flag to auto-created directory
#[test]
fn test_compress_with_auto_create_directory() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Create non-existent destination directory
    let dest_dir = temp_dir.path().join("output");
    assert!(!dest_dir.exists());

    // Compress with -C to non-existent directory (should auto-create)
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    // Verify directory was created
    assert!(
        dir_exists(&dest_dir),
        "Destination directory should be created"
    );

    // Verify compressed file is in the destination
    let compressed_file = dest_dir.join("test.txt.gz");
    assert!(
        file_exists(&compressed_file),
        "Compressed file should exist in destination"
    );
}

/// Test compression with nested non-existent directory path
#[test]
fn test_compress_with_nested_directory_creation() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    // Create nested non-existent path
    let dest_dir = temp_dir.path().join("deep").join("nested").join("path");
    assert!(!dest_dir.exists());

    // Compress to nested path
    jcz_command()
        .arg("-c")
        .arg("bzip2")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    // Verify nested directories were created
    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("test.txt.bz2")));
}

/// Test cross-device compression with -C flag
#[test]
fn test_cross_device_compress_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Create a separate temp directory to simulate cross-device
    let dest_temp_dir = TempDir::new().unwrap();
    let dest_dir = dest_temp_dir.path().join("output");
    fs::create_dir(&dest_dir).unwrap();

    // Compress with -C to different location
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    // Verify file was moved/copied correctly
    assert!(file_exists(&dest_dir.join("test.txt.gz")));
}

/// Test compression of multiple files with -C flag
#[test]
fn test_compress_multiple_files_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let files = create_test_files(
        temp_dir.path(),
        &[
            ("file1.txt", TEST_DATA_SMALL),
            ("file2.txt", TEST_DATA_MEDIUM),
            ("file3.txt", TEST_DATA_BINARY),
        ],
    );

    let dest_dir = temp_dir.path().join("compressed");

    // Compress multiple files to destination
    jcz_command()
        .arg("-c")
        .arg("xz")
        .arg("-C")
        .arg(&dest_dir)
        .args(&files)
        .assert()
        .success();

    // Verify all compressed files are in destination
    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("file1.txt.xz")));
    assert!(file_exists(&dest_dir.join("file2.txt.xz")));
    assert!(file_exists(&dest_dir.join("file3.txt.xz")));

    // Verify originals still exist
    assert!(file_exists(&files[0]));
    assert!(file_exists(&files[1]));
    assert!(file_exists(&files[2]));
}

/// Test compound format compression with -C flag
#[test]
fn test_compress_compound_format_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    let dest_dir = temp_dir.path().join("archives");

    // Test tar.gz
    jcz_command()
        .arg("-c")
        .arg("tgz")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("test.txt.tar.gz")));

    // Test tar.bz2
    jcz_command()
        .arg("-c")
        .arg("tbz2")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(file_exists(&dest_dir.join("test.txt.tar.bz2")));

    // Test tar.xz
    jcz_command()
        .arg("-c")
        .arg("txz")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(file_exists(&dest_dir.join("test.txt.tar.xz")));
}

/// Test compression with -C to existing directory
#[test]
fn test_compress_to_existing_directory() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Create destination directory
    let dest_dir = temp_dir.path().join("existing");
    fs::create_dir(&dest_dir).unwrap();

    // Compress to existing directory
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(file_exists(&dest_dir.join("test.txt.gz")));
}

/// Test compression with -C and different compression levels
#[test]
fn test_compress_with_c_flag_and_levels() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    let dest_dir = temp_dir.path().join("compressed");

    // Test with level 1
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-l")
        .arg("1")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(dir_exists(&dest_dir));
    let compressed_file = dest_dir.join("test.txt.gz");
    assert!(file_exists(&compressed_file));

    // Remove compressed file
    fs::remove_file(&compressed_file).unwrap();

    // Test with level 9
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-l")
        .arg("9")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(file_exists(&compressed_file));
}

/// Test compression preserves original file
#[test]
fn test_compress_with_c_flag_preserves_original() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let original_content = read_file(&test_file);

    let dest_dir = temp_dir.path().join("output");

    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    // Verify original still exists with same content
    assert!(file_exists(&test_file), "Original file should be preserved");
    assert_eq!(
        read_file(&test_file),
        original_content,
        "Original content should be unchanged"
    );

    // Verify compressed file exists in destination
    assert!(file_exists(&dest_dir.join("test.txt.gz")));
}

/// Test compression with -C flag and timestamp option
#[test]
fn test_compress_with_c_flag_and_timestamp() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    let dest_dir = temp_dir.path().join("timestamped");

    // Compress with timestamp option
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-t")
        .arg("1") // Date timestamp
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(dir_exists(&dest_dir));

    // Check that a .gz file exists (we can't predict exact timestamp filename)
    let entries: Vec<_> = fs::read_dir(&dest_dir)
        .unwrap()
        .filter_map(|e| e.ok())
        .filter(|e| e.path().extension().map_or(false, |ext| ext == "gz"))
        .collect();

    assert_eq!(entries.len(), 1, "Should have exactly one .gz file");
}

/// Test concurrent compression to different destinations
#[test]
fn test_concurrent_compress_with_different_destinations() {
    let temp_dir = TempDir::new().unwrap();

    let file1 = create_test_file(temp_dir.path(), "test1.txt", TEST_DATA_SMALL);
    let file2 = create_test_file(temp_dir.path(), "test2.txt", TEST_DATA_MEDIUM);

    let dest1 = temp_dir.path().join("dest1");
    let dest2 = temp_dir.path().join("dest2");

    // Compress first file to dest1
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&dest1)
        .arg(&file1)
        .assert()
        .success();

    // Compress second file to dest2
    jcz_command()
        .arg("-c")
        .arg("bzip2")
        .arg("-C")
        .arg(&dest2)
        .arg(&file2)
        .assert()
        .success();

    // Verify both destinations have correct files
    assert!(file_exists(&dest1.join("test1.txt.gz")));
    assert!(file_exists(&dest2.join("test2.txt.bz2")));
}

/// Test compression with -C flag handles binary data correctly
#[test]
fn test_compress_binary_data_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "binary.dat", TEST_DATA_BINARY);

    let dest_dir = temp_dir.path().join("binary_output");

    jcz_command()
        .arg("-c")
        .arg("xz")
        .arg("-C")
        .arg(&dest_dir)
        .arg(&test_file)
        .assert()
        .success();

    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("binary.dat.xz")));

    // Verify original binary data is preserved
    assert_eq!(read_file(&test_file), TEST_DATA_BINARY);
}
