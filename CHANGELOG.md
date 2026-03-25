# Changelog

## v0.8.0 (2026-03-26)

### Features
- LogFileOutput option to tools for enhanced logging.
- Various updates to workflows, argument parsing, build, and utility files.

### Breaking
- Switched to use external goext package (see github.com/roemer/goext)
- Removed own execution runner (execr) -> use the new goext CmdRunner
- Refactored all NxTool methods to accept settings as pointers for consistency.
