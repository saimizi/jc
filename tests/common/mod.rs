use assert_cmd::Command;
use std::fs;
use std::path::{Path, PathBuf};
use tempfile::TempDir;

/// Helper to create a test file with specified content
pub fn create_test_file(dir: &Path, name: &str, content: &[u8]) -> PathBuf {
    let file_path = dir.join(name);
    fs::write(&file_path, content).expect("Failed to create test file");
    file_path
}

/// Helper to create multiple test files
#[allow(dead_code)]
pub fn create_test_files(dir: &Path, files: &[(&str, &[u8])]) -> Vec<PathBuf> {
    files
        .iter()
        .map(|(name, content)| create_test_file(dir, name, content))
        .collect()
}

/// Helper to create a jcz command
pub fn jcz_command() -> Command {
    #[allow(deprecated)]
    Command::cargo_bin("jcz").expect("Failed to find jcz binary")
}

/// Helper to verify a file exists
pub fn file_exists(path: &Path) -> bool {
    path.exists() && path.is_file()
}

/// Helper to verify directory exists
#[allow(dead_code)]
pub fn dir_exists(path: &Path) -> bool {
    path.exists() && path.is_dir()
}

/// Helper to get file size
#[allow(dead_code)]
pub fn file_size(path: &Path) -> u64 {
    fs::metadata(path)
        .expect("Failed to get file metadata")
        .len()
}

/// Helper to read file content
#[allow(dead_code)]
pub fn read_file(path: &Path) -> Vec<u8> {
    fs::read(path).expect("Failed to read file")
}

/// Helper to verify file has specific extension
#[allow(dead_code)]
pub fn has_extension(path: &Path, ext: &str) -> bool {
    path.extension()
        .and_then(|e| e.to_str())
        .map(|e| e == ext)
        .unwrap_or(false)
}

/// Helper to decompress a file using system tools and verify content
#[allow(dead_code)]
pub fn verify_decompressed_content(compressed_path: &Path, expected_content: &[u8]) -> bool {
    let temp_dir = TempDir::new().expect("Failed to create temp dir");
    let _output_path = temp_dir.path().join("output");

    // Determine decompression command based on extension
    let result = if compressed_path.to_str().unwrap().ends_with(".gz") {
        std::process::Command::new("gzip")
            .args(["-dc", compressed_path.to_str().unwrap()])
            .output()
    } else if compressed_path.to_str().unwrap().ends_with(".bz2") {
        std::process::Command::new("bzip2")
            .args(["-dc", compressed_path.to_str().unwrap()])
            .output()
    } else if compressed_path.to_str().unwrap().ends_with(".xz") {
        std::process::Command::new("xz")
            .args(["-dc", compressed_path.to_str().unwrap()])
            .output()
    } else {
        return false;
    };

    match result {
        Ok(output) if output.status.success() => output.stdout == expected_content,
        _ => false,
    }
}

/// Helper to create a test directory structure
#[allow(dead_code)]
pub fn create_test_dir_structure(base: &Path, structure: &[&str]) -> Vec<PathBuf> {
    let mut paths = Vec::new();
    for path_str in structure {
        let path = base.join(path_str);
        if path_str.ends_with('/') {
            fs::create_dir_all(&path).expect("Failed to create directory");
        } else {
            if let Some(parent) = path.parent() {
                fs::create_dir_all(parent).expect("Failed to create parent directory");
            }
            fs::write(&path, format!("Content of {}", path_str)).expect("Failed to create file");
            paths.push(path);
        }
    }
    paths
}

/// Test data content
pub const TEST_DATA_SMALL: &[u8] = b"Hello, World! This is a test file for compression.";
#[allow(dead_code)]
pub const TEST_DATA_MEDIUM: &[u8] = b"Lorem ipsum dolor sit amet, consectetur adipiscing elit. \
    Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \
    Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris. \
    This content is repeated to make it more compressible. \
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. \
    Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.";
#[allow(dead_code)]
pub const TEST_DATA_BINARY: &[u8] = &[0u8, 1, 2, 3, 4, 5, 255, 254, 253, 252, 251, 250];
