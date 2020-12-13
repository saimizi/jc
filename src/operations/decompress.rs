use rayon::prelude::*;
use std::path::PathBuf;

use crate::compressors::{create_compressor, detect_format};
use crate::core::config::CompressionConfig;
use crate::core::error::{JcError, JcResult};
use crate::utils::{debug, error, info, remove_file_silent};

/// Decompress a single file, handling compound formats
pub fn decompress_file(input: &PathBuf, config: &CompressionConfig) -> JcResult<PathBuf> {
    let mut current_file = input.clone();
    let mut temp_files = Vec::new();

    // Iteratively decompress until no more compression detected
    loop {
        let format = detect_format(&current_file).ok_or_else(|| {
            JcError::InvalidExtension(
                current_file.clone(),
                "supported compression format".to_string(),
            )
        })?;

        debug!(
            "Detected format: {:?} for {}",
            format,
            current_file.display()
        );

        let compressor = create_compressor(format);
        let output = compressor.decompress(&current_file, config)?;

        // If this was an intermediate file, mark for cleanup
        if current_file != *input {
            temp_files.push(current_file.clone());
        }

        current_file = output;

        // Check if output has another compression layer
        if detect_format(&current_file).is_none() {
            break;
        }
    }

    // Clean up temporary intermediate files
    for temp in temp_files {
        if let Err(e) = remove_file_silent(&temp) {
            debug!("Failed to remove temp file {}: {}", temp.display(), e);
        }
    }

    Ok(current_file)
}

/// Decompress multiple files concurrently
pub fn decompress_files(inputs: Vec<PathBuf>, config: CompressionConfig) -> Vec<JcResult<PathBuf>> {
    info!("Decompressing {} files", inputs.len());

    inputs
        .par_iter()
        .map(|input| match decompress_file(input, &config) {
            Ok(output) => Ok(output),
            Err(e) => {
                error!("Failed to decompress {}: {}", input.display(), e);
                Err(e)
            }
        })
        .collect()
}
