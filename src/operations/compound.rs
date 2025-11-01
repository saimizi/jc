use rayon::prelude::*;
use std::path::PathBuf;

use crate::compressors::create_compressor;
use crate::core::config::CompressionConfig;
use crate::core::error::JcResult;
use crate::core::types::CompoundFormat;
use crate::utils::{debug, info, remove_file_silent};

/// Compress file(s) with compound format (TAR + secondary compression)
pub fn compress_compound(
    input: &PathBuf,
    format: CompoundFormat,
    config: &CompressionConfig,
) -> JcResult<PathBuf> {
    info!(
        "Compressing {} with compound format: {}",
        input.display(),
        format.extension()
    );

    // Step 1: Create TAR archive
    let tar_compressor = create_compressor(format.primary());
    let tar_config = CompressionConfig {
        level: 0, // TAR doesn't use compression level
        timestamp: config.timestamp,
        move_to: None, // Don't move intermediate file
        show_output_size: false,
        force: config.force,
    };

    let tar_output = tar_compressor.compress(input, &tar_config)?;
    debug!("Created intermediate TAR: {}", tar_output.display());

    // Step 2: Compress TAR with secondary compressor
    let secondary_compressor = create_compressor(format.secondary());
    let secondary_output = secondary_compressor.compress(&tar_output, config)?;

    // Step 3: Remove intermediate TAR file
    if let Err(e) = remove_file_silent(&tar_output) {
        debug!("Failed to remove intermediate TAR: {}", e);
    }

    info!("Created compound archive: {}", secondary_output.display());
    Ok(secondary_output)
}

/// Compress multiple files with compound format
pub fn compress_compound_batch(
    inputs: Vec<PathBuf>,
    format: CompoundFormat,
    config: CompressionConfig,
) -> Vec<JcResult<PathBuf>> {
    inputs
        .par_iter()
        .map(|input| compress_compound(input, format, &config))
        .collect()
}
