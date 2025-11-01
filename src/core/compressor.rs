use std::path::{Path, PathBuf};

use crate::core::config::CompressionConfig;
use crate::core::error::JcResult;

/// Common interface for all compression/decompression implementations
#[allow(dead_code)]
pub trait Compressor: Send + Sync {
    /// Get the name of this compressor
    fn name(&self) -> &'static str;

    /// Get the file extension used by this compressor (e.g., "gz", "bz2")
    fn extension(&self) -> &'static str;

    /// Compress a single file
    fn compress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf>;

    /// Decompress a single file
    fn decompress(&self, input: &Path, config: &CompressionConfig) -> JcResult<PathBuf>;

    /// Check if this compressor supports compression levels
    fn supports_levels(&self) -> bool;

    /// Validate compression level for this compressor
    fn validate_level(&self, level: u8) -> bool;

    /// Get default compression level
    fn default_level(&self) -> u8;
}

/// Extended trait for compressors that support multi-file operations
#[allow(dead_code)]
pub trait MultiFileCompressor: Compressor {
    /// Compress multiple files into a single archive
    fn compress_multi(
        &self,
        inputs: &[PathBuf],
        output_name: &str,
        config: &CompressionConfig,
    ) -> JcResult<PathBuf>;
}
