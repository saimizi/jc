mod common;

use common::*;
use std::fs;
use tempfile::TempDir;

/// Test decompression with automatic directory creation (-C flag)
#[test]
fn test_decompress_with_auto_create_directory() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // First compress the file
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");
    assert!(file_exists(&compressed_file));

    // Remove the original file
    fs::remove_file(&test_file).unwrap();

    // Create a non-existent destination directory path
    let dest_dir = temp_dir.path().join("output");
    assert!(!dest_dir.exists(), "Destination should not exist initially");

    // Decompress to non-existent directory (should auto-create)
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-C")
        .arg(&dest_dir)
        .arg("-f")
        .assert()
        .success();

    // Verify directory was created
    assert!(
        dir_exists(&dest_dir),
        "Destination directory should be created"
    );

    // Verify file was decompressed
    let decompressed_file = dest_dir.join("test.txt");
    assert!(
        file_exists(&decompressed_file),
        "Decompressed file should exist"
    );
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

/// Test decompression with nested non-existent directory path
#[test]
fn test_decompress_with_nested_directory_creation() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress the file
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");
    fs::remove_file(&test_file).unwrap();

    // Create nested non-existent path
    let dest_dir = temp_dir.path().join("deep").join("nested").join("path");
    assert!(!dest_dir.exists());

    // Decompress to nested path
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-C")
        .arg(&dest_dir)
        .arg("-f")
        .assert()
        .success();

    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("test.txt")));
}

/// Test decompression with --force flag (skip overwrite prompts)
#[test]
fn test_decompress_with_force_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress the file
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");

    // Decompress first time
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-f")
        .assert()
        .success();

    // Modify the decompressed file
    let decompressed_file = temp_dir.path().join("test.txt");
    fs::write(&decompressed_file, b"Modified content").unwrap();
    assert_eq!(read_file(&decompressed_file), b"Modified content");

    // Decompress again with --force (should overwrite without prompt)
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-f")
        .assert()
        .success();

    // Verify it was overwritten
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

/// Test decompression of compound format (.tar.gz) with single file
#[test]
fn test_decompress_compound_format_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    // Compress to .tar.gz
    jcz_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(file_exists(&compressed_file));

    // Remove original and decompress
    fs::remove_file(&test_file).unwrap();

    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-f")
        .assert()
        .success();

    // Verify decompressed file
    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_MEDIUM);
}

/// Test decompression with -C flag to different directory
#[test]
fn test_decompress_to_different_directory_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress the file
    jcz_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(file_exists(&compressed_file));

    // Create destination directory
    let dest_dir = temp_dir.path().join("extracted");
    fs::create_dir(&dest_dir).unwrap();

    // Remove original file
    fs::remove_file(&test_file).unwrap();

    // Decompress to different directory with -C flag
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-C")
        .arg(&dest_dir)
        .arg("-f")
        .assert()
        .success();

    // Verify file was extracted to dest_dir
    assert!(file_exists(&dest_dir.join("test.txt")));
    assert_eq!(read_file(&dest_dir.join("test.txt")), TEST_DATA_SMALL);
}

/// Test concurrent decompression without file conflicts
#[test]
fn test_concurrent_decompress_no_conflicts() {
    let temp_dir = TempDir::new().unwrap();

    // Create and compress multiple files
    let file1 = create_test_file(temp_dir.path(), "test1.txt", TEST_DATA_SMALL);
    let file2 = create_test_file(temp_dir.path(), "test2.txt", TEST_DATA_MEDIUM);
    let file3 = create_test_file(temp_dir.path(), "test3.txt", TEST_DATA_BINARY);

    jcz_command()
        .arg("-c")
        .arg("tgz")
        .args(&[&file1, &file2, &file3])
        .assert()
        .success();

    let compressed1 = temp_dir.path().join("test1.txt.tar.gz");
    let compressed2 = temp_dir.path().join("test2.txt.tar.gz");
    let compressed3 = temp_dir.path().join("test3.txt.tar.gz");

    // Remove originals
    fs::remove_file(&file1).unwrap();
    fs::remove_file(&file2).unwrap();
    fs::remove_file(&file3).unwrap();

    // Decompress all files concurrently (should not conflict due to temp directory isolation)
    jcz_command()
        .arg("-d")
        .arg(&compressed1)
        .arg(&compressed2)
        .arg(&compressed3)
        .arg("-f")
        .assert()
        .success();

    // Verify all were decompressed successfully
    assert!(file_exists(&temp_dir.path().join("test1.txt")));
    assert!(file_exists(&temp_dir.path().join("test2.txt")));
    assert!(file_exists(&temp_dir.path().join("test3.txt")));
    assert_eq!(
        read_file(&temp_dir.path().join("test1.txt")),
        TEST_DATA_SMALL
    );
    assert_eq!(
        read_file(&temp_dir.path().join("test2.txt")),
        TEST_DATA_MEDIUM
    );
    assert_eq!(
        read_file(&temp_dir.path().join("test3.txt")),
        TEST_DATA_BINARY
    );
}

/// Test decompression with compound format (.tar.bz2)
#[test]
fn test_decompress_tar_bz2() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    // Compress to .tar.bz2
    jcz_command()
        .arg("-c")
        .arg("tbz2")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.bz2");
    assert!(file_exists(&compressed_file));

    fs::remove_file(&test_file).unwrap();

    // Decompress
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-f")
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_MEDIUM);
}

/// Test decompression with compound format (.tar.xz)
#[test]
fn test_decompress_tar_xz() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    // Compress to .tar.xz
    jcz_command()
        .arg("-c")
        .arg("txz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.xz");
    assert!(file_exists(&compressed_file));

    fs::remove_file(&test_file).unwrap();

    // Decompress
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-f")
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_MEDIUM);
}

/// Test decompression of different formats to same directory
#[test]
fn test_decompress_mixed_formats() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    // Compress with different formats
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    jcz_command()
        .arg("-c")
        .arg("bzip2")
        .arg(&test_file)
        .assert()
        .success();

    jcz_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let gz_file = temp_dir.path().join("test.txt.gz");
    let bz2_file = temp_dir.path().join("test.txt.bz2");
    let xz_file = temp_dir.path().join("test.txt.xz");

    // Remove original
    fs::remove_file(&test_file).unwrap();

    // Decompress gz first
    jcz_command()
        .arg("-d")
        .arg(&gz_file)
        .arg("-f")
        .assert()
        .success();

    assert!(file_exists(&test_file));
    assert_eq!(read_file(&test_file), TEST_DATA_MEDIUM);

    // Decompress bz2 (should overwrite with force flag)
    jcz_command()
        .arg("-d")
        .arg(&bz2_file)
        .arg("-f")
        .assert()
        .success();

    assert_eq!(read_file(&test_file), TEST_DATA_MEDIUM);

    // Decompress xz (should overwrite with force flag)
    jcz_command()
        .arg("-d")
        .arg(&xz_file)
        .arg("-f")
        .assert()
        .success();

    assert_eq!(read_file(&test_file), TEST_DATA_MEDIUM);
}

/// Test cross-device move with -C flag (simulated by using temp directory)
#[test]
fn test_cross_device_move_with_c_flag() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress the file
    jcz_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");

    // Create a separate temp directory for destination
    let dest_temp_dir = TempDir::new().unwrap();
    let dest_dir = dest_temp_dir.path().join("output");

    // Decompress with -C to different location
    jcz_command()
        .arg("-d")
        .arg(&compressed_file)
        .arg("-C")
        .arg(&dest_dir)
        .arg("-f")
        .assert()
        .success();

    // Verify file was moved/copied correctly
    assert!(dir_exists(&dest_dir));
    assert!(file_exists(&dest_dir.join("test.txt")));
    assert_eq!(read_file(&dest_dir.join("test.txt")), TEST_DATA_SMALL);
}

/// Test decompression preserves file content exactly (including binary data)
#[test]
fn test_decompress_preserves_binary_content() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "binary.dat", TEST_DATA_BINARY);

    // Compress with each algorithm
    for format in &["gzip", "bzip2", "xz"] {
        jcz_command()
            .arg("-c")
            .arg(format)
            .arg(&test_file)
            .assert()
            .success();

        let ext = match *format {
            "gzip" => "gz",
            "bzip2" => "bz2",
            "xz" => "xz",
            _ => unreachable!(),
        };

        let compressed_file = temp_dir.path().join(format!("binary.dat.{}", ext));
        assert!(file_exists(&compressed_file));

        // Decompress
        fs::remove_file(&test_file).unwrap();

        jcz_command()
            .arg("-d")
            .arg(&compressed_file)
            .arg("-f")
            .assert()
            .success();

        // Verify binary content is preserved
        assert!(file_exists(&test_file));
        assert_eq!(read_file(&test_file), TEST_DATA_BINARY);

        // Clean up for next iteration
        fs::remove_file(&compressed_file).unwrap();
    }
}
