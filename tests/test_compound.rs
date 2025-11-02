mod common;

use common::*;
use tempfile::TempDir;

// TGZ Tests (TAR + GZIP)

#[test]
fn test_tgz_compress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(file_exists(&compressed_file), "TGZ file should exist");
    assert!(file_exists(&test_file), "Original file should be preserved");
}

#[test]
fn test_tgz_compress_with_level_1() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg("-l")
        .arg("1")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_tgz_compress_with_level_9() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg("-l")
        .arg("9")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_tgz_decompress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress
    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Decompress
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

#[test]
fn test_tgz_default_command() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // tgz is the default, so omitting -c should use tgz
    jc_command().arg(&test_file).assert().success();

    let compressed_file = temp_dir.path().join("test.txt.tar.gz");
    assert!(
        file_exists(&compressed_file),
        "TGZ file should exist as default"
    );
}

// TBZ2 Tests (TAR + BZIP2)

#[test]
fn test_tbz2_compress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.bz2");
    assert!(file_exists(&compressed_file), "TBZ2 file should exist");
    assert!(file_exists(&test_file), "Original file should be preserved");
}

#[test]
fn test_tbz2_compress_with_level_1() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg("-l")
        .arg("1")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.bz2");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_tbz2_compress_with_level_9() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg("-l")
        .arg("9")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.bz2");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_tbz2_decompress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress
    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.bz2");

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Decompress
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

// TXZ Tests (TAR + XZ)

#[test]
fn test_txz_compress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.xz");
    assert!(file_exists(&compressed_file), "TXZ file should exist");
    assert!(file_exists(&test_file), "Original file should be preserved");
}

#[test]
fn test_txz_compress_with_level_1() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg("-l")
        .arg("1")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.xz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_txz_compress_with_level_9() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_MEDIUM);

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg("-l")
        .arg("9")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.xz");
    assert!(file_exists(&compressed_file));
}

#[test]
fn test_txz_decompress_single_file() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    // Compress
    jc_command()
        .arg("-c")
        .arg("txz")
        .arg(&test_file)
        .assert()
        .success();

    let compressed_file = temp_dir.path().join("test.txt.tar.xz");

    // Remove original
    std::fs::remove_file(&test_file).unwrap();

    // Decompress
    jc_command()
        .arg("-d")
        .arg(&compressed_file)
        .assert()
        .success();

    let decompressed_file = temp_dir.path().join("test.txt");
    assert!(file_exists(&decompressed_file));
    assert_eq!(read_file(&decompressed_file), TEST_DATA_SMALL);
}

// Compound format verification tests

#[test]
fn test_tgz_verify_is_gzip_compressed() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();

    let tgz_file = temp_dir.path().join("test.txt.tar.gz");

    // Verify it's a gzip file by checking magic number
    let content = read_file(&tgz_file);
    assert!(
        content.len() >= 2 && content[0] == 0x1f && content[1] == 0x8b,
        "TGZ file should have gzip magic number"
    );
}

#[test]
fn test_tbz2_verify_is_bzip2_compressed() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg(&test_file)
        .assert()
        .success();

    let tbz2_file = temp_dir.path().join("test.txt.tar.bz2");

    // Verify it's a bzip2 file by checking magic number
    let content = read_file(&tbz2_file);
    assert!(
        content.len() >= 3 && content[0] == b'B' && content[1] == b'Z' && content[2] == b'h',
        "TBZ2 file should have bzip2 magic number"
    );
}

#[test]
fn test_txz_verify_is_xz_compressed() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg(&test_file)
        .assert()
        .success();

    let txz_file = temp_dir.path().join("test.txt.tar.xz");

    // Verify it's an xz file by checking magic number
    let content = read_file(&txz_file);
    assert!(
        content.len() >= 6
            && content[0] == 0xfd
            && content[1] == 0x37
            && content[2] == 0x7a
            && content[3] == 0x58
            && content[4] == 0x5a
            && content[5] == 0x00,
        "TXZ file should have xz magic number"
    );
}

#[test]
fn test_compound_formats_preserve_originals() {
    let temp_dir = TempDir::new().unwrap();
    let test_file = create_test_file(temp_dir.path(), "test.txt", TEST_DATA_SMALL);
    let original_content = read_file(&test_file);

    // Test TGZ
    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg(&test_file)
        .assert()
        .success();
    assert_eq!(read_file(&test_file), original_content);

    // Test TBZ2
    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg(&test_file)
        .assert()
        .success();
    assert_eq!(read_file(&test_file), original_content);

    // Test TXZ
    jc_command()
        .arg("-c")
        .arg("txz")
        .arg(&test_file)
        .assert()
        .success();
    assert_eq!(read_file(&test_file), original_content);
}

// Multi-file archive tests (-A flag)

#[test]
fn test_tgz_compress_multiple_files_with_archive_flag() {
    let temp_dir = TempDir::new().unwrap();
    let file1 = create_test_file(temp_dir.path(), "file1.txt", b"Content 1");
    let file2 = create_test_file(temp_dir.path(), "file2.txt", b"Content 2");
    let file3 = create_test_file(temp_dir.path(), "file3.txt", b"Content 3");

    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg("-A")
        .arg("archive")
        .arg(&file1)
        .arg(&file2)
        .arg(&file3)
        .arg("-C")
        .arg(temp_dir.path())
        .assert()
        .success();

    let archive_file = temp_dir.path().join("archive.tar.gz");
    assert!(file_exists(&archive_file), "TGZ archive should exist");

    // Verify it's a valid gzip file
    let content = read_file(&archive_file);
    assert!(
        content.len() >= 2 && content[0] == 0x1f && content[1] == 0x8b,
        "TGZ archive should have gzip magic number"
    );
}

#[test]
fn test_tbz2_compress_multiple_files_with_archive_flag() {
    let temp_dir = TempDir::new().unwrap();
    let file1 = create_test_file(temp_dir.path(), "file1.txt", b"Content 1");
    let file2 = create_test_file(temp_dir.path(), "file2.txt", b"Content 2");
    let file3 = create_test_file(temp_dir.path(), "file3.txt", b"Content 3");

    jc_command()
        .arg("-c")
        .arg("tbz2")
        .arg("-A")
        .arg("archive")
        .arg(&file1)
        .arg(&file2)
        .arg(&file3)
        .arg("-C")
        .arg(temp_dir.path())
        .assert()
        .success();

    let archive_file = temp_dir.path().join("archive.tar.bz2");
    assert!(file_exists(&archive_file), "TBZ2 archive should exist");

    // Verify it's a valid bzip2 file
    let content = read_file(&archive_file);
    assert!(
        content.len() >= 3 && content[0] == b'B' && content[1] == b'Z' && content[2] == b'h',
        "TBZ2 archive should have bzip2 magic number"
    );
}

#[test]
fn test_txz_compress_multiple_files_with_archive_flag() {
    let temp_dir = TempDir::new().unwrap();
    let file1 = create_test_file(temp_dir.path(), "file1.txt", b"Content 1");
    let file2 = create_test_file(temp_dir.path(), "file2.txt", b"Content 2");
    let file3 = create_test_file(temp_dir.path(), "file3.txt", b"Content 3");

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg("-A")
        .arg("archive")
        .arg(&file1)
        .arg(&file2)
        .arg(&file3)
        .arg("-C")
        .arg(temp_dir.path())
        .assert()
        .success();

    let archive_file = temp_dir.path().join("archive.tar.xz");
    assert!(file_exists(&archive_file), "TXZ archive should exist");

    // Verify it's a valid xz file
    let content = read_file(&archive_file);
    assert!(
        content.len() >= 6
            && content[0] == 0xfd
            && content[1] == 0x37
            && content[2] == 0x7a
            && content[3] == 0x58
            && content[4] == 0x5a
            && content[5] == 0x00,
        "TXZ archive should have xz magic number"
    );
}

#[test]
fn test_multiple_files_archive_with_destination_directory() {
    let temp_dir = TempDir::new().unwrap();
    let source_dir = temp_dir.path().join("source");
    let dest_dir = temp_dir.path().join("dest");
    std::fs::create_dir(&source_dir).unwrap();
    std::fs::create_dir(&dest_dir).unwrap();

    let file1 = create_test_file(&source_dir, "a.txt", b"Content A");
    let file2 = create_test_file(&source_dir, "b.txt", b"Content B");
    let file3 = create_test_file(&source_dir, "c.txt", b"Content C");

    jc_command()
        .arg("-c")
        .arg("txz")
        .arg("-A")
        .arg("myarchive")
        .arg(&file1)
        .arg(&file2)
        .arg(&file3)
        .arg("-C")
        .arg(&dest_dir)
        .assert()
        .success();

    let archive_file = dest_dir.join("myarchive.tar.xz");
    assert!(
        file_exists(&archive_file),
        "Archive should be created in destination directory"
    );
}

#[test]
fn test_multiple_files_archive_can_be_extracted() {
    let temp_dir = TempDir::new().unwrap();
    let file1 = create_test_file(temp_dir.path(), "file1.txt", b"Content 1");
    let file2 = create_test_file(temp_dir.path(), "file2.txt", b"Content 2");
    let file3 = create_test_file(temp_dir.path(), "file3.txt", b"Content 3");

    // Create archive
    jc_command()
        .arg("-c")
        .arg("tgz")
        .arg("-A")
        .arg("test_archive")
        .arg(&file1)
        .arg(&file2)
        .arg(&file3)
        .arg("-C")
        .arg(temp_dir.path())
        .assert()
        .success();

    let archive_file = temp_dir.path().join("test_archive.tar.gz");
    assert!(file_exists(&archive_file));

    // Remove original files
    std::fs::remove_file(&file1).unwrap();
    std::fs::remove_file(&file2).unwrap();
    std::fs::remove_file(&file3).unwrap();

    // Extract using system tar to verify archive integrity
    std::process::Command::new("tar")
        .arg("-xzf")
        .arg(&archive_file)
        .arg("-C")
        .arg(temp_dir.path())
        .output()
        .expect("Failed to extract archive");

    // Verify extracted files
    assert!(file_exists(&temp_dir.path().join("file1.txt")));
    assert!(file_exists(&temp_dir.path().join("file2.txt")));
    assert!(file_exists(&temp_dir.path().join("file3.txt")));
    assert_eq!(read_file(&temp_dir.path().join("file1.txt")), b"Content 1");
    assert_eq!(read_file(&temp_dir.path().join("file2.txt")), b"Content 2");
    assert_eq!(read_file(&temp_dir.path().join("file3.txt")), b"Content 3");
}
