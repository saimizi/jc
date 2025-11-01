use std::path::{Path, PathBuf};
use std::process::Command;

use crate::core::compressor::{Compressor, MultiFileCompressor};
use crate::core::config::CompressionConfig;
use crate::core::error::{JcError, JcResult};
use crate::utils::{copy_to_dir, debug, generate_output_filename, info, move_file_if_needed};

/// TAR archiver implementation
#[derive(Debug, Clone)]
pub struct TarCompressor;

impl TarCompressor {
    pub fn new() -> Self {
        Self
    }
}

impl Compressor for TarCompressor {
    fn name(&self) -> &'static str {
        "tar"
    }

    fn extension(&self) -> &'static str {
        "tar"
    }

    fn compress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf> {
        if !input.exists() {
            return Err(JcError::FileNotFound(input.to_path_buf()));
        }

        let output_path = generate_output_filename(input, "tar", config.timestamp)?;
        info!(
            "Creating TAR archive {} from {}",
            output_path.display(),
            input.display()
        );

        // Build tar command - if input has no parent, use current directory
        let mut cmd = Command::new("tar");

        if let Some(parent) = input.parent().filter(|p| !p.as_os_str().is_empty()) {
            let basename = input
                .file_name()
                .ok_or_else(|| JcError::Other("Invalid filename".to_string()))?;
            cmd.arg("-C").arg(parent);
            cmd.arg("-cf").arg(&output_path).arg(basename);
        } else {
            // No parent or empty parent, just use the input path directly
            cmd.arg("-cf").arg(&output_path).arg(input);
        }

        debug!("Executing: {:?}", cmd);

        let output = cmd
            .output()
            .map_err(|e| JcError::Other(format!("Failed to execute tar: {}", e)))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::CompressionFailed {
                tool: "tar".to_string(),
                stderr: stderr.to_string(),
            });
        }

        // Move to destination if specified
        let final_path = move_file_if_needed(&output_path, &config.move_to)?;

        info!("Created TAR archive: {}", final_path.display());
        Ok(final_path)
    }

    fn decompress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf> {
        if !input.to_string_lossy().ends_with(".tar") {
            return Err(JcError::InvalidExtension(
                input.to_path_buf(),
                "tar".to_string(),
            ));
        }

        debug!("Extracting TAR archive {}", input.display());

        let parent = input
            .parent()
            .filter(|p| !p.as_os_str().is_empty())
            .unwrap_or_else(|| Path::new("."));

        let mut cmd = Command::new("tar");
        cmd.arg("-x").arg("-C").arg(parent).arg("-f").arg(input);

        let output = cmd
            .output()
            .map_err(|e| JcError::Other(format!("Failed to execute tar: {}", e)))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::DecompressionFailed {
                tool: "tar".to_string(),
                stderr: stderr.to_string(),
            });
        }

        // Output is the filename without .tar extension
        let output_path = input.with_extension("");

        // Move to destination if specified
        let final_path = move_file_if_needed(&output_path, &config.move_to)?;

        info!("Extracted TAR archive to: {}", final_path.display());
        Ok(final_path)
    }

    fn supports_levels(&self) -> bool {
        false // TAR doesn't support compression levels
    }

    fn validate_level(&self, _level: u8) -> bool {
        true // Always valid (no-op)
    }

    fn default_level(&self) -> u8 {
        0
    }
}

impl TarCompressor {
    /// Decompress in a specific working directory
    pub fn decompress_in_dir(
        &self,
        input: &Path,
        working_dir: &Path,
        _config: &CompressionConfig,
    ) -> JcResult<PathBuf> {
        if !input.to_string_lossy().ends_with(".tar") {
            return Err(JcError::InvalidExtension(
                input.to_path_buf(),
                "tar".to_string(),
            ));
        }

        debug!(
            "Extracting TAR archive {} in working dir {}",
            input.display(),
            working_dir.display()
        );

        // Copy input file to working directory
        let work_input = copy_to_dir(input, working_dir)?;

        // Extract in working directory
        let mut cmd = Command::new("tar");
        cmd.arg("-x")
            .arg("-C")
            .arg(working_dir)
            .arg("-f")
            .arg(&work_input);

        let output = cmd
            .output()
            .map_err(|e| JcError::Other(format!("Failed to execute tar: {}", e)))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::DecompressionFailed {
                tool: "tar".to_string(),
                stderr: stderr.to_string(),
            });
        }

        // TAR extracts files, so we need to find what was extracted
        // List the directory to find extracted content
        use std::fs;
        let entries: Vec<_> = fs::read_dir(working_dir)
            .map_err(|e| JcError::Io(e))?
            .filter_map(|e| e.ok())
            .filter(|e| e.path() != work_input) // Exclude the tar file itself
            .collect();

        // If we found exactly one entry, use that
        if entries.len() == 1 {
            let extracted_path = entries[0].path();
            debug!("Extracted to: {}", extracted_path.display());
            return Ok(extracted_path);
        }

        // Multiple files extracted - check if there's a common parent directory
        // or a directory with the same base name as the tar file
        let tar_base_name = work_input
            .file_stem()
            .and_then(|s| s.to_str())
            .unwrap_or("");

        // Check if there's a directory with the tar's base name
        for entry in &entries {
            let path = entry.path();
            if path.is_dir() {
                if let Some(dir_name) = path.file_name().and_then(|s| s.to_str()) {
                    if dir_name == tar_base_name {
                        debug!("Extracted to directory: {}", path.display());
                        return Ok(path);
                    }
                }
            }
        }

        // If we have multiple files but no matching directory, just return the working directory
        // This happens when tar extracts multiple loose files
        if !entries.is_empty() {
            // Remove the tar file itself to avoid copying it to the destination
            use std::fs;
            let _ = fs::remove_file(&work_input);

            debug!(
                "Extracted {} files to: {}",
                entries.len(),
                working_dir.display()
            );
            return Ok(working_dir.to_path_buf());
        }

        // Fallback: assume filename without .tar extension (original behavior)
        let output_path = work_input.with_extension("");
        debug!("Extracted to (fallback): {}", output_path.display());
        Ok(output_path)
    }
}

impl MultiFileCompressor for TarCompressor {
    fn compress_multi(
        &self,
        inputs: &[PathBuf],
        output_name: &str,
        config: &CompressionConfig,
    ) -> JcResult<PathBuf> {
        if inputs.is_empty() {
            return Err(JcError::NoInputFiles);
        }

        let mut output_path = PathBuf::from(output_name);
        if !output_path.extension().map_or(false, |e| e == "tar") {
            output_path.set_extension("tar");
        }

        info!("Creating multi-file TAR archive: {}", output_path.display());

        let mut cmd = Command::new("tar");
        cmd.arg("-cf").arg(&output_path);

        for input in inputs {
            let basename = input
                .file_name()
                .ok_or_else(|| JcError::Other("Invalid filename".to_string()))?;
            cmd.arg(basename);
        }

        debug!("Executing: {:?}", cmd);

        let output = cmd
            .output()
            .map_err(|e| JcError::Other(format!("Failed to execute tar: {}", e)))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::CompressionFailed {
                tool: "tar".to_string(),
                stderr: stderr.to_string(),
            });
        }

        // Move to destination if specified
        let final_path = move_file_if_needed(&output_path, &config.move_to)?;

        Ok(final_path)
    }
}
