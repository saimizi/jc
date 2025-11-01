mod common;

use common::*;
use tempfile::TempDir;

#[test]
fn test_xz_compress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    assert!(
        file_exists(&compressed_file),
        "Compressed file should exist"
    );
    assert!(file_exists(&test_file), "Original file should be preserved");
}

#[test]
fn test_xz_compress_multiple_files() {
    let temp_dir = TempDir::new().unwrap();
    let files = create_test_files(
        temp_dir.path(),
        &[
            ("file1.txt", TEST_DATA_SMALL),
            ("file2.txt", TEST_DATA_MEDIUM),
        ],
    );

    jc_command()
        .arg("-c")
        .arg("xz")
        .args(&files)
        .assert()
        .success();

    assert!(file_exists(&temp_dir.path().join("file1.txt.xz")));
    assert!(file_exists(&temp_dir.path().join("file2.txt.xz")));
    assert!(file_exists(&files[0]), "Original files should be preserved");
    assert!(file_exists(&files[1]), "Original files should be preserved");
}

#[test]
fn test_xz_compress_with_level_1() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg("-l")
        .arg("1")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_xz_compress_with_level_9() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg("-l")
        .arg("9")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_xz_compress_with_default_level() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_xz_decompress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // First compress
    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    assert!(file_exists(&compressed_file));

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Now decompress
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(
        file_exists(&decompressed_file),
        "Decompressed file should exist"
    );
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

#[test]
fn test_xz_decompress_multiple_files() {
    let temp_dir = TempDir::new().unwrap();
    let files = create_test_files(
        temp_dir.path(),
        &[
            ("file1.txt", TEST_DATA_SMALL),
            ("file2.txt", TEST_DATA_MEDIUM),
        ],
    );

    // Compress both files
    jc_command()
        .arg("-c")
        .arg("xz")
        .args(&files)
        .assert()
        .success();

    let xz1 = temp_dir.path().join("file1.txt.xz");
    let xz2 = temp_dir.path().join("file2.txt.xz");

    // Remove originals
    std::fs::remove_file(&files[0]).unwrap();
    std::fs::remove_file(&files[1]).unwrap();

    // Decompress both
    jc_command()
        .arg("-d")
        .arg(&xz1)
        .arg(&xz2)
        .assert()
        .success();

    assert!(file_exists(&files[0]));
    assert!(file_exists(&files[1]));
    assert_eq!(read_file(&files[0]), TEST_DATA_SMALL);
    assert_eq!(read_file(&files[1]), TEST_DATA_MEDIUM);
}

#[test]
fn test_xz_compress_binary_data() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "binary.dat", TEST_DATA_BINARY);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("binary.dat.xz");
    assert!(file_exists(&compressed_file));

    // Verify decompression works
    std::fs::remove_file(&test_file).unwrap();
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .success();

    assert_eq!(
        read_file(&temp_dir.path().join("binary.dat")),
        TEST_DATA_BINARY
    );
}

#[test]
fn test_xz_verify_compression_reduces_size() {
    let temp_dir = TempDir::new().unwrap();

    // Create a highly compressible file (repeated content)
    let compressible_content = "X".repeat(10000);
    let test_file = create_test_file(temp_dir.path(), "test.txt", compressible_content.as_bytes());
    let original_size = file_size(&test_file);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.xz");
    let compressed_size = file_size(&compressed_file);

    assert!(
        compressed_size < original_size,
        "Compressed file should be smaller than original"
    );
}

#[test]
fn test_xz_compress_preserves_original() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let original_content = read_file(&test_file);

    jc_command()
        .arg("-c")
        .arg("xz")
        .arg(&test_file)
        .assert()
        .success();

    assert!(file_exists(&test_file), "Original file should still exist");
    assert_eq!(
        read_file(&test_file),
        original_content,
        "Original file content should be unchanged"
    );
}
