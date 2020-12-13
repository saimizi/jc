pub mod cli;
pub mod compressors;
pub mod core;
pub mod operations;
pub mod utils;

// Re-export commonly used types for library users
pub use core::{
    CollectionConfig, CollectionMode, CompressionConfig, CompressionFormat, Compressor, JcError,
    JcResult, TimestampOption,
};

pub use operations::{
    collect_and_compress, compress_compound, compress_file, compress_files, decompress_file,
    decompress_files,
};
