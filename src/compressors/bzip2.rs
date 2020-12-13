use std::fs::File;
use std::io::{BufWriter, Write};
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};

use crate::core::compressor::Compressor;
use crate::core::config::CompressionConfig;
use crate::core::error::{JcError, JcResult};
use crate::utils::{debug, generate_output_filename, info, move_file_if_needed};

/// BZIP2 compressor implementation
#[derive(Debug, Clone)]
pub struct Bzip2Compressor;

impl Bzip2Compressor {
    pub fn new() -> Self {
        Self
    }

    fn validate_input(&self, path: &Path) -> JcResult<()> {
        if !path.exists() {
            return Err(JcError::FileNotFound(path.to_path_buf()));
        }

        if path.is_dir() {
            return Err(JcError::NotAFile(path.to_path_buf()));
        }

        Ok(())
    }
}

impl Compressor for Bzip2Compressor {
    fn name(&self) -> &'static str {
        "bzip2"
    }

    fn extension(&self) -> &'static str {
        "bz2"
    }

    fn compress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf> {
        self.validate_input(input)?;

        let output_path = generate_output_filename(input, "bz2", config.timestamp)?;
        info!(
            "Compressing {} to {} with bzip2",
            input.display(),
            output_path.display()
        );
        debug!("Compression level: {}", config.level);

        let output_file = File::create(&output_path)?;
        let mut writer = BufWriter::new(output_file);

        let mut cmd = Command::new("bzip2");
        cmd.arg(format!("-{}", config.level))
            .arg("--keep")
            .arg("--stdout")
            .arg(input)
            .stdout(Stdio::piped())
            .stderr(Stdio::piped());

        debug!("Executing: {:?}", cmd);

        let mut child = cmd
            .spawn()
            .map_err(|e| JcError::Other(format!("Failed to spawn bzip2: {}", e)))?;

        if let Some(mut stdout) = child.stdout.take() {
            std::io::copy(&mut stdout, &mut writer)?;
        }

        writer.flush()?;

        let output = child.wait_with_output()?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::CompressionFailed {
                tool: "bzip2".to_string(),
                stderr: stderr.to_string(),
            });
        }

        let final_path = move_file_if_needed(&output_path, &config.move_to)?;

        info!("Compressed file: {}", final_path.display());
        Ok(final_path)
    }

    fn decompress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf> {
        if !input.to_string_lossy().ends_with(".bz2") {
            return Err(JcError::InvalidExtension(
                input.to_path_buf(),
                "bz2".to_string(),
            ));
        }

        debug!("Decompressing {} with bzip2", input.display());

        let mut cmd = Command::new("bzip2");
        cmd.arg("-d").arg("-k").arg(input);

        let output = cmd
            .output()
            .map_err(|e| JcError::Other(format!("Failed to execute bzip2: {}", e)))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            return Err(JcError::DecompressionFailed {
                tool: "bzip2".to_string(),
                stderr: stderr.to_string(),
            });
        }

        let output_path = input.with_extension("");
        let final_path = move_file_if_needed(&output_path, &config.move_to)?;

        info!("Decompressed file: {}", final_path.display());
        Ok(final_path)
    }

    fn supports_levels(&self) -> bool {
        true
    }

    fn validate_level(&self, level: u8) -> bool {
        (1..=9).contains(&level)
    }

    fn default_level(&self) -> u8 {
        6
    }
}
