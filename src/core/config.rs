use std::path::PathBuf;

/// Timestamp formatting options
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TimestampOption {
    None,        // 0: No timestamp
    Date,        // 1: YYYYMMDD
    DateTime,    // 2: YYYYMMDD_HHMMSS
    Nanoseconds, // 3: Nanoseconds only
}

impl TimestampOption {
    pub fn from_u8(value: u8) -> Option<Self> {
        match value {
            0 => Some(TimestampOption::None),
            1 => Some(TimestampOption::Date),
            2 => Some(TimestampOption::DateTime),
            3 => Some(TimestampOption::Nanoseconds),
            _ => None,
        }
    }
}

/// Configuration for compression/decompression operations
#[derive(Debug, Clone)]
pub struct CompressionConfig {
    /// Compression level (0-9, meaning varies by algorithm)
    pub level: u8,

    /// Timestamp option for output filenames
    pub timestamp: TimestampOption,

    /// Destination directory for output files
    pub move_to: Option<PathBuf>,

    /// Show output file size (future feature)
    #[allow(dead_code)]
    pub show_output_size: bool,

    /// Force overwrite without prompting
    pub force: bool,
}

impl Default for CompressionConfig {
    fn default() -> Self {
        Self {
            level: 6,
            timestamp: TimestampOption::None,
            move_to: None,
            show_output_size: false,
            force: false,
        }
    }
}

impl CompressionConfig {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn with_level(mut self, level: u8) -> Self {
        self.level = level;
        self
    }

    pub fn with_timestamp(mut self, timestamp: TimestampOption) -> Self {
        self.timestamp = timestamp;
        self
    }

    pub fn with_move_to(mut self, path: PathBuf) -> Self {
        self.move_to = Some(path);
        self
    }

    pub fn with_force(mut self, force: bool) -> Self {
        self.force = force;
        self
    }
}

/// Collection operation mode
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum CollectionMode {
    /// Include parent directory in archive (-a flag)
    WithParent,

    /// Archive files without parent directory wrapper (-A flag)
    Flat,
}

/// Configuration for collection operations (multi-file archives)
#[derive(Debug, Clone)]
pub struct CollectionConfig {
    /// Base configuration
    pub base: CompressionConfig,

    /// Package/archive name
    pub package_name: String,

    /// Collection mode
    pub mode: CollectionMode,
}
