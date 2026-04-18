# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Support for passing album parameter to Last.fm when scrobbling tracks (both CLI and TUI modes)
- `--album` flag to track command in CLI mode
- Optional album field in TUI mode for track scrobbling
- Album validation - when user provides `--album`, validates it's a valid release for the track via Last.fm

### Changed
- Output now displays album name in parentheses for both track and album scrobbles
- Artist, track, and album names are now corrected using Last.fm metadata for proper casing
- Track command always fetches metadata to get corrected names from Last.fm

## [0.1.1] - 2024-05-15

### Added
- Profile management (add, list, set, get, delete, open)
- Track scrobbling (CLI and TUI modes)
- Album scrobbling (CLI and TUI modes)
- TUI mode with interactive forms for artist/track/album and timestamp selection
- Last.fm API integration (search, scrobble, history)
- Session key storage in system keychain
- `--end-now` flag to calculate start time from track/album duration
- `--date` and `--timestamp` flags for custom scrobble times
- `--dryrun` flag to preview scrobbles without submitting
- `-p/--profile` flag to specify profile
- History command to review recent scrobbles

### Changed
- Reformatted output for consistent display
- Track timestamps are calculated relative to duration rather than all at once for albums

### Fixed
- Track command parsing error
- Date and time fields in album TUI
- Form validation for empty fields

### Removed
- Debug flag and logging in release builds