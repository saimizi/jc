use std::fs;
use std::io;
use std::path::{Path, PathBuf};

use crate::core::config::TimestampOption;
use crate::core::error::{JcError, JcResult};
use crate::utils::timestamp::generate_timestamp;

/// Generate output filename with optional timestamp
pub fn generate_output_filename(
    input: &Path,
    extension: &str,
    timestamp_opt: TimestampOption,
) -> JcResult<PathBuf> {
    let mut filename = input.as_os_str().to_string_lossy().to_string();

    // Remove trailing slash if present
    if filename.ends_with('/') {
        filename.pop();
    }

    // Add timestamp if requested
    if timestamp_opt != TimestampOption::None {
        let ts = generate_timestamp(timestamp_opt);
        filename.push('_');
        filename.push_str(&ts);
    }

    // Add extension
    filename.push('.');
    filename.push_str(extension);

    Ok(PathBuf::from(filename))
}

/// Move file to destination directory if specified
pub fn move_file_if_needed(source: &Path, move_to: &Option<PathBuf>) -> JcResult<PathBuf> {
    if let Some(dest_dir) = move_to {
        move_file(source, dest_dir)
    } else {
        Ok(source.to_path_buf())
    }
}

/// Move file to destination directory
pub fn move_file(source: &Path, dest_dir: &Path) -> JcResult<PathBuf> {
    // Validate destination is a directory
    if !dest_dir.is_dir() {
        return Err(JcError::NotADirectory(dest_dir.to_path_buf()));
    }

    let filename = source
        .file_name()
        .ok_or_else(|| JcError::Other("Invalid source filename".to_string()))?;

    let dest_path = dest_dir.join(filename);

    fs::rename(source, &dest_path)?;

    Ok(dest_path)
}

/// Recursively copy file or directory
pub fn copy_recursive(src: &Path, dst: &Path) -> io::Result<()> {
    if src.is_dir() {
        fs::create_dir_all(dst)?;
        for entry in fs::read_dir(src)? {
            let entry = entry?;
            let src_path = entry.path();
            let dst_path = dst.join(entry.file_name());
            copy_recursive(&src_path, &dst_path)?;
        }
    } else {
        fs::copy(src, dst)?;
    }
    Ok(())
}

/// Remove file, ignoring errors
pub fn remove_file_silent(path: &Path) -> io::Result<()> {
    match fs::remove_file(path) {
        Ok(()) => Ok(()),
        Err(e) if e.kind() == io::ErrorKind::NotFound => Ok(()),
        Err(e) => Err(e),
    }
}

/// Create temporary directory with prefix
pub fn create_temp_dir(prefix: &str) -> JcResult<PathBuf> {
    use std::time::{SystemTime, UNIX_EPOCH};

    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_nanos();

    let dir_name = format!("{}{:x}", prefix, timestamp);
    let temp_path = PathBuf::from(dir_name);

    fs::create_dir(&temp_path).map_err(|e| JcError::TempDirFailed(e.to_string()))?;

    Ok(temp_path)
}
