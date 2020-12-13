use chrono::Local;

use crate::core::config::TimestampOption;

/// Generate timestamp string based on option
pub fn generate_timestamp(option: TimestampOption) -> String {
    let now = Local::now();

    match option {
        TimestampOption::None => String::new(),
        TimestampOption::Date => now.format("%Y%m%d").to_string(),
        TimestampOption::DateTime => now.format("%Y%m%d_%H%M%S").to_string(),
        TimestampOption::Nanoseconds => now.timestamp_subsec_nanos().to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_timestamp_none() {
        let ts = generate_timestamp(TimestampOption::None);
        assert_eq!(ts, "");
    }

    #[test]
    fn test_timestamp_date() {
        let ts = generate_timestamp(TimestampOption::Date);
        assert_eq!(ts.len(), 8); // YYYYMMDD
        assert!(ts.chars().all(|c| c.is_ascii_digit()));
    }

    #[test]
    fn test_timestamp_datetime() {
        let ts = generate_timestamp(TimestampOption::DateTime);
        assert_eq!(ts.len(), 15); // YYYYMMDD_HHMMSS
        assert!(ts.contains('_'));
    }
}
