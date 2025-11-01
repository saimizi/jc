mod common;

use common::*;
use tempfile::TempDir;

// Invalid Option Tests

#[test]
fn test_invalid_compression_command() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("invalid")
        .arg(&test_file)
        .assert()
        .failure();
}

#[test]
fn test_invalid_compression_level_zero() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-l")
        .arg("0")
        .arg(&test_file)
        .assert()
        .failure();
}

#[test]
fn test_invalid_compression_level_too_high() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-l")
        .arg("10")
        .arg(&test_file)
        .assert()
        .failure();
}

#[test]
fn test_invalid_timestamp_option_negative() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-t")
        .arg("-1")
        .arg(&test_file)
        .assert()
        .failure();
}

#[test]
fn test_invalid_timestamp_option_too_high() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-t")
        .arg("4")
        .arg(&test_file)
        .assert()
        .failure();
}

// Missing File Tests

#[test]
fn test_nonexistent_file() {
    let temp_dir = TempDir::new().unwrap();
    let nonexistent = temp_dir.path().join("does_not_exist.txt");

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&nonexistent)
        .assert()
        .failure();
}

#[test]
fn test_decompress_nonexistent_file() {
    let temp_dir = TempDir::new().unwrap();
    let nonexistent = temp_dir.path().join("does_not_exist.gz");

    jc_command().arg("-d").arg(&nonexistent).assert().failure();
}

#[test]
fn test_multiple_files_with_one_missing() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "exists.txt", TEST_DATA_SMALL);
    let nonexistent = temp_dir.path().join("does_not_exist.txt");

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .arg(&nonexistent)
        .assert()
        .failure();
}

// No Input Files Tests

#[test]
fn test_no_input_files() {
    jc_command().arg("-c").arg("gzip").assert().failure();
}

#[test]
fn test_decompress_no_input_files() {
    jc_command().arg("-d").assert().failure();
}

// Invalid Directory Tests

#[test]
fn test_move_to_nonexistent_directory() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let nonexistent_dir = temp_dir.path().join("does_not_exist");

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&nonexistent_dir)
        .arg(&test_file)
        .assert()
        .failure();
}

#[test]
fn test_move_to_file_instead_of_directory() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let not_a_dir = create_test_file(temp_dir.path(), "not_a_dir.txt", b"content");

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg("-C")
        .arg(&not_a_dir)
        .arg(&test_file)
        .assert()
        .failure();
}

// Collection Option Conflicts

#[test]
fn test_collect_without_archive_name() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Using -a flag requires an argument
    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg("-a")
        .arg(&test_file)
        .assert()
        .failure();
}

// Invalid Decompression Tests

#[test]
fn test_decompress_invalid_gzip_file() {
    let temp_dir = TempDir::new().unwrap();
    let fake_gz = create_test_file(temp_dir.path(), "fake.gz", b"not a gzip file");

    jc_command().arg("-d").arg(&fake_gz).assert().failure();
}

#[test]
fn test_decompress_corrupted_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // First create a valid compressed file
    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");

    // Corrupt the file by truncating it
    let mut content = read_file(&compressed_file);
    content.truncate(content.len() / 2);
    std::fs::write(&compressed_file, content).unwrap();

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Try to decompress corrupted file
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .failure();
}

// Directory as Input Tests

#[test]
fn test_compress_directory_without_collection() {
    let temp_dir = TempDir::new().unwrap();
    let test_dir = temp_dir.path().join("testdir");
    std::fs::create_dir(&test_dir).unwrap();
    create_test_file(&test_dir, "file.txt", TEST_DATA_SMALL);

    // Trying to compress a directory without collection should fail or handle appropriately
    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_dir)
        .assert()
        .failure();
}

// Edge Case Tests

#[test]
fn test_empty_file() {
    let temp_dir = TempDir::new().unwrap();
    let empty_file = create_test_file(temp_dir.path(), "empty.txt", b"");

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&empty_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("empty.txt.gz");
    assert!(
        file_exists(&compressed_file),
        "Empty file should be compressed"
    );
}

#[test]
fn test_file_with_special_characters_in_name() {
    let temp_dir = TempDir::new().unwrap();
    let special_file = create_test_file(
        temp_dir.path(),
        "test file with spaces.txt",
        TEST_DATA_SMALL,
    );

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&special_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test file with spaces.txt.gz");
    assert!(
        file_exists(&compressed_file),
        "File with spaces should be compressed"
    );
}

#[test]
fn test_file_with_multiple_dots() {
    let temp_dir = TempDir::new().unwrap();
    let dotted_file = create_test_file(temp_dir.path(), "file.tar.backup.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&dotted_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("file.tar.backup.txt.gz");
    assert!(
        file_exists(&compressed_file),
        "File with multiple dots should be compressed"
    );
}

#[test]
fn test_already_compressed_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress once
    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");

    // Try to compress the .gz file (should work, creating .gz.gz)
    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&compressed_file)
        .assert()
        .success();

    let double_compressed = temp_dir.path().join("test.txt.gz.gz");
    assert!(
        file_exists(&double_compressed),
        "Should allow compressing already compressed file"
    );
}

// Decompression with wrong extension

#[test]
fn test_decompress_file_with_wrong_extension() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress with gzip
    jc_command()
        .arg("-c")
        .arg("gzip")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.gz");

    // Rename to wrong extension
    let wrong_ext = temp_dir.path().join("test.txt.bz2");
    std::fs::rename(&compressed_file, &wrong_ext).unwrap();

    // Try to decompress - jc decompresses based on extension, so this should fail
    jc_command().arg("-d").arg(&wrong_ext).assert().failure();
}
