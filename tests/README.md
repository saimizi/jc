# Integration Tests

This directory contains integration tests for the `jc` (Just Compress) command-line tool.

## Prerequisites

The integration tests require the following system tools to be installed:
- `gzip`
- `bzip2`
- `xz`
- `tar`

These are typically pre-installed on most Linux systems and GitHub runners.

## Running Tests

### Run all integration tests:
```bash
cargo test --test '*'
```

### Run a specific test file:
```bash
cargo test --test test_gzip
cargo test --test test_bzip2
cargo test --test test_xz
cargo test --test test_tar
cargo test --test test_compound
cargo test --test test_options
cargo test --test test_errors
```

### Run a specific test:
```bash
cargo test --test test_gzip test_gzip_compress_single_file
```

### Run with output:
```bash
cargo test --test '*' -- --nocapture
```

### Run tests in parallel (default):
```bash
cargo test --test '*'
```

### Run tests sequentially:
```bash
cargo test --test '*' -- --test-threads=1
```

## Test Organization

Tests are organized by compression format and feature:

- **test_gzip.rs** - GZIP compression and decompression tests
- **test_bzip2.rs** - BZIP2 compression and decompression tests
- **test_xz.rs** - XZ compression and decompression tests
- **test_tar.rs** - TAR archive tests
- **test_compound.rs** - Compound format tests (TGZ, TBZ2, TXZ)
- **test_options.rs** - Cross-cutting options tests (timestamp, move-to, collection)
- **test_errors.rs** - Error handling and edge case tests
- **common/mod.rs** - Shared test utilities and helper functions

## CI/CD Integration

These tests are designed to run on x86 GitHub runners without special hardware requirements.

To integrate with GitHub Actions, add the following to your workflow:

```yaml
- name: Run integration tests
  run: cargo test --test '*'
```

## Test Coverage

The test suite includes **112 tests** covering:
- All compression formats (gzip, bzip2, xz, tar, tgz, tbz2, txz)
- All compression levels (1-9)
- Timestamp options (0-3)
- Move-to directory option (-C) *
- Collection modes (-a, -A) *
- Compress and decompress operations
- Multiple file handling
- Binary data handling
- Error conditions and edge cases

\* **Note on skipped tests:** Some tests for the move-to directory (`-C`) and collection (`-a`, `-A`) options are currently skipped due to cross-filesystem rename issues when using `tempfile` crate. These tests are marked with `#[ignore]` and can be run individually with `--ignored` flag if needed on a single filesystem.
