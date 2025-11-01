use env_logger::{Builder, Env};
use log::LevelFilter;
use std::sync::OnceLock;

static LOGGER_INIT: OnceLock<()> = OnceLock::new();

/// Initialize the logging system based on JCDBG environment variable
pub fn init_logger() {
    LOGGER_INIT.get_or_init(|| {
        let env = Env::default()
            .filter_or("JCDBG", "info")
            .write_style("JCDBG_STYLE");

        let mut builder = Builder::from_env(env);

        // Map JCDBG values to log levels
        let level = std::env::var("JCDBG")
            .ok()
            .and_then(|val| match val.to_lowercase().as_str() {
                "error" => Some(LevelFilter::Error),
                "warn" => Some(LevelFilter::Warn),
                "info" => Some(LevelFilter::Info),
                "debug" => Some(LevelFilter::Debug),
                _ => None,
            })
            .unwrap_or(LevelFilter::Info);

        builder
            .filter_level(level)
            .format_module_path(false)
            .format_target(false)
            .init();
    });
}

// Re-export log macros for convenience
pub use log::{debug, error, info};
