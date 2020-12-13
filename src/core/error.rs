use std::fmt;
use std::io;
use std::path::PathBuf;

/// Result type for JC operations
pub type JcResult<T> = Result<T, JcError>;

/// Comprehensive error type for JC operations
#[derive(Debug)]
pub enum JcError {
    /// File not found
    FileNotFound(PathBuf),

    /// Path is not a file (e.g., directory when file expected)
    NotAFile(PathBuf),

    /// Path is not a directory
    NotADirectory(PathBuf),

    /// Invalid file extension for operation
    InvalidExtension(PathBuf, String),

    /// Invalid compression level for algorithm
    InvalidCompressionLevel { algorithm: String, level: u8 },

    /// Invalid timestamp option
    InvalidTimestampOption(u8),

    /// Invalid compression command
    InvalidCommand(String),

    /// Duplicate basenames in collection
    DuplicateBasenames(Vec<String>),

    /// Archive/package name already exists
    NameExists(String),

    /// Move-to directory error
    MoveToError(String),

    /// Compression tool execution failed
    CompressionFailed { tool: String, stderr: String },

    /// Decompression tool execution failed
    DecompressionFailed { tool: String, stderr: String },

    /// I/O error
    Io(io::Error),

    /// Symbolic link resolution failed
    SymlinkResolution(PathBuf),

    /// Temporary directory creation failed
    TempDirFailed(String),

    /// No input files provided
    NoInputFiles,

    /// Generic error with message
    Other(String),
}

impl fmt::Display for JcError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            JcError::FileNotFound(path) => {
                write!(f, "File not found: {}", path.display())
            }
            JcError::NotAFile(path) => {
                write!(f, "{} is not a file", path.display())
            }
            JcError::NotADirectory(path) => {
                write!(f, "{} is not a directory", path.display())
            }
            JcError::InvalidExtension(path, expected) => {
                write!(
                    f,
                    "{} has invalid extension, expected: {}",
                    path.display(),
                    expected
                )
            }
            JcError::InvalidCompressionLevel { algorithm, level } => {
                write!(f, "Invalid compression level {} for {}", level, algorithm)
            }
            JcError::InvalidTimestampOption(opt) => {
                write!(f, "Invalid timestamp option: {}", opt)
            }
            JcError::InvalidCommand(cmd) => {
                write!(f, "Invalid compression command: {}", cmd)
            }
            JcError::DuplicateBasenames(names) => {
                write!(f, "Duplicate basenames in collection: {}", names.join(", "))
            }
            JcError::NameExists(name) => {
                write!(
                    f,
                    "{} already exists and cannot be used as package name",
                    name
                )
            }
            JcError::MoveToError(msg) => {
                write!(f, "Move-to directory error: {}", msg)
            }
            JcError::CompressionFailed { tool, stderr } => {
                write!(f, "{} compression failed: {}", tool, stderr)
            }
            JcError::DecompressionFailed { tool, stderr } => {
                write!(f, "{} decompression failed: {}", tool, stderr)
            }
            JcError::Io(err) => {
                write!(f, "I/O error: {}", err)
            }
            JcError::SymlinkResolution(path) => {
                write!(f, "Failed to resolve symbolic link: {}", path.display())
            }
            JcError::TempDirFailed(msg) => {
                write!(f, "Temporary directory creation failed: {}", msg)
            }
            JcError::NoInputFiles => {
                write!(f, "No input files provided")
            }
            JcError::Other(msg) => {
                write!(f, "{}", msg)
            }
        }
    }
}

impl std::error::Error for JcError {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        match self {
            JcError::Io(err) => Some(err),
            _ => None,
        }
    }
}

impl From<io::Error> for JcError {
    fn from(err: io::Error) -> Self {
        JcError::Io(err)
    }
}
