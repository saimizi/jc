use rayon::prelude::*;
use std::path::PathBuf;

use crate::compressors::create_compressor;
use crate::core::config::CompressionConfig;
use crate::core::error::{JcError, JcResult};
use crate::core::types::CompressionFormat;
use crate::utils::{error, info};

/// Compress a single file
pub fn compress_file(
    input: &PathBuf,
    format: CompressionFormat,
    config: &CompressionConfig,
) -> JcResult<PathBuf> {
    let compressor = create_compressor(format);

    // Validate compression level if supported
    if compressor.supports_levels() && !compressor.validate_level(config.level) {
        return Err(JcError::InvalidCompressionLevel {
            algorithm: compressor.name().to_string(),
            level: config.level,
        });
    }

    compressor.compress(input, config)
}

/// Compress multiple files concurrently
pub fn compress_files(
    inputs: Vec<PathBuf>,
    format: CompressionFormat,
    config: CompressionConfig,
) -> Vec<JcResult<PathBuf>> {
    info!("Compressing {} files with {}", inputs.len(), format.name());

    // Use rayon for parallel processing
    inputs
        .par_iter()
        .map(|input| match compress_file(input, format, &config) {
            Ok(output) => Ok(output),
            Err(e) => {
                error!("Failed to compress {}: {}", input.display(), e);
                Err(e)
            }
        })
        .collect()
}
