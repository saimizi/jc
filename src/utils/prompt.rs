use std::io::{self, Write};
use std::path::Path;

use crate::core::error::{JcError, JcResult};

/// Prompt user whether to overwrite an existing file
/// Returns true if user confirms, false otherwise
pub fn prompt_overwrite(file_path: &Path) -> JcResult<bool> {
    print!(
        "File '{}' already exists. Overwrite? (y/n): ",
        file_path.display()
    );
    io::stdout().flush().map_err(|e| JcError::Io(e))?;

    let mut input = String::new();
    io::stdin()
        .read_line(&mut input)
        .map_err(|e| JcError::Io(e))?;

    let response = input.trim().to_lowercase();
    Ok(response == "y" || response == "yes")
}
