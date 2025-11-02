mod common;

use common::*;
use std::process::Command;
use tempfile::TempDir;

#[test]
fn test_tar_archive_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jcz_command()
        .arg("-c")
        .arg("tar")
        .arg(&test_file)
        .assert()
        .success();

    let archive_file = temp_dir.path().join("test.txt.tar");
    assert!(file_exists(&archive_file), "Archive file should exist");
    assert!(file_exists(&test_file), "Original file should be preserved");
}

#[test]
fn test_tar_archive_multiple_files() {
    let temp_dir = TempDir::new().unwrap();
    let files = create_test_files(
        temp_dir.path(),
        &[
            ("file1.txt", TEST_DATA_SMALL),
            ("file2.txt", TEST_DATA_MEDIUM),
        ],
    );

    jcz_command()
        .arg("-c")
        .arg("tar")
        .args(&files)
        .assert()
        .success();

    // TAR with multiple files creates separate archives
    assert!(file_exists(&temp_dir.path().join("file1.txt.tar")));
    assert!(file_exists(&temp_dir.path().join("file2.txt.tar")));
    assert!(file_exists(&files[0]), "Original files should be preserved");
    assert!(file_exists(&files[1]), "Original files should be preserved");
}

#[test]
fn test_tar_extract_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Create archive
    jcz_command()
        .arg("-c")
        .arg("tar")
        .arg(&test_file)
        .assert()
        .success();

    let archive_file = temp_dir.path().join("test.txt.tar");
    assert!(file_exists(&archive_file));

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Extract
    jcz_command()
        .arg("-d")
        .arg(&archive_file)
        .assert()
        .success();

    let extracted_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&extracted_file), "Extracted file should exist");
    assert_eq!(read_file(&extracted_file), TEST_DATA_SMALL);
}

#[test]
fn test_tar_extract_multiple_files() {
    let temp_dir = TempDir::new().unwrap();
    let files = create_test_files(
        temp_dir.path(),
        &[
            ("file1.txt", TEST_DATA_SMALL),
            ("file2.txt", TEST_DATA_MEDIUM),
        ],
    );

    // Create archives
    jcz_command()
        .arg("-c")
        .arg("tar")
        .args(&files)
        .assert()
        .success();

    let tar1 = temp_dir.path().join("file1.txt.tar");
    let tar2 = temp_dir.path().join("file2.txt.tar");

    // Remove originals
    std::fs::remove_file(&files[0]).unwrap();
    std::fs::remove_file(&files[1]).unwrap();

    // Extract both
    jcz_command()
        .arg("-d")
        .arg(&tar1)
        .arg(&tar2)
        .assert()
        .success();

    assert!(file_exists(&files[0]));
    assert!(file_exists(&files[1]));
    assert_eq!(read_file(&files[0]), TEST_DATA_SMALL);
    assert_eq!(read_file(&files[1]), TEST_DATA_MEDIUM);
}

#[test]
fn test_tar_archive_preserves_original() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let original_content = read_file(&test_file);

    jcz_command()
        .arg("-c")
        .arg("tar")
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

#[test]
fn test_tar_archive_binary_data() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "binary.dat", TEST_DATA_BINARY);

    jcz_command()
        .arg("-c")
        .arg("tar")
        .arg(&test_file)
        .assert()
        .success();

    let archive_file = temp_dir.path().join("binary.dat.tar");
    assert!(file_exists(&archive_file));

    // Verify extraction works
    std::fs::remove_file(&test_file).unwrap();
    jcz_command()
        .arg("-d")
        .arg(&archive_file)
        .assert()
        .success();

    assert_eq!(
        read_file(&temp_dir.path().join("binary.dat")),
        TEST_DATA_BINARY
    );
}

#[test]
fn test_tar_verify_archive_contains_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let file_name = test_file.file_name().unwrap().to_str().unwrap();

    jcz_command()
        .arg("-c")
        .arg("tar")
        .arg(&test_file)
        .assert()
        .success();

    let archive_file = temp_dir.path().join("test.txt.tar");

    // List contents of tar to verify file is included
    let output = Command::new("tar")
        .arg("-tf")
        .arg(&archive_file)
        .output()
        .expect("Failed to list tar contents");

    assert!(output.status.success());
    let contents = String::from_utf8_lossy(&output.stdout);
    assert!(
        contents.contains(file_name),
        "Archive should contain the original file"
    );
}

#[test]
fn test_tar_with_compression_level_ignored() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // TAR doesn't use compression level, but command should not fail
    jcz_command()
        .arg("-c")
        .arg("tar")
        .arg("-l")
        .arg("9")
        .arg(&test_file)
        .assert()
        .success();

    let archive_file = temp_dir.path().join("test.txt.tar");
    assert!(file_exists(&archive_file));
}
