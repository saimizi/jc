use std::collections::{HashMap, HashSet};
use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;

use crate::core::error::{JcError, JcResult};
use crate::core::types::InputFile;
use crate::utils::logger::debug;

/// Validate and process input files
pub fn validate_input_files(paths: Vec<PathBuf>) -> JcResult<Vec<InputFile>> {
    if paths.is_empty() {
        return Err(JcError::NoInputFiles);
    }

    let mut validated = Vec::new();
    let mut seen_paths = HashSet::new();

    for path in paths {
        // Check if file exists
        let metadata = fs::metadata(&path).map_err(|_| JcError::FileNotFound(path.clone()))?;

        // Resolve symbolic links
        let (real_path, was_symlink) = if metadata.file_type().is_symlink() {
            debug!("{} is a symbolic link, resolving", path.display());
            let real = resolve_symlink(&path)?;
            (real, true)
        } else {
            (path.clone(), false)
        };

        // Check for duplicates
        if !seen_paths.insert(real_path.clone()) {
            debug!("Skipping duplicate path: {}", real_path.display());
            continue;
        }

        // Get basename
        let basename = real_path
            .file_name()
            .and_then(|n| n.to_str())
            .ok_or_else(|| JcError::Other("Invalid filename".to_string()))?
            .to_string();

        validated.push(InputFile {
            original_path: path,
            real_path,
            basename,
            was_symlink,
        });
    }

    Ok(validated)
}

/// Check if there are duplicate basenames
pub fn check_duplicate_basenames(files: &[InputFile]) -> Option<Vec<String>> {
    let mut basename_counts: HashMap<&str, usize> = HashMap::new();

    for file in files {
        *basename_counts.entry(&file.basename).or_insert(0) += 1;
    }

    let duplicates: Vec<String> = basename_counts
        .iter()
        .filter(|(_, &count)| count > 1)
        .map(|(name, _)| name.to_string())
        .collect();

    if duplicates.is_empty() {
        None
    } else {
        Some(duplicates)
    }
}

/// Resolve symbolic link to real path
fn resolve_symlink(path: &Path) -> JcResult<PathBuf> {
    let output = Command::new("readlink")
        .arg("-f")
        .arg(path)
        .output()
        .map_err(|_| JcError::SymlinkResolution(path.to_path_buf()))?;

    if !output.status.success() {
        return Err(JcError::SymlinkResolution(path.to_path_buf()));
    }

    let real_path = String::from_utf8_lossy(&output.stdout);
    let real_path = real_path.trim();

    Ok(PathBuf::from(real_path))
}

/// Validate destination directory
pub fn validate_move_to(path: &Path) -> JcResult<()> {
    if !path.exists() {
        return Err(JcError::MoveToError("Directory does not exist".to_string()));
    }

    if !path.is_dir() {
        return Err(JcError::NotADirectory(path.to_path_buf()));
    }

    // Check if writable by attempting to create a test file
    let test_file = path.join(".jc_write_test");
    match fs::File::create(&test_file) {
        Ok(_) => {
            let _ = fs::remove_file(&test_file);
            Ok(())
        }
        Err(_) => Err(JcError::MoveToError(
            "Directory is not writable".to_string(),
        )),
    }
}
