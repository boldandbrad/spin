# Contributing

## Reporting Bugs

Open an [issue](https://github.com/boldandbrad/spin/issues) with details.

## Pull Requests

1. Fork the repo
2. Make your changes
3. Run `go build ./... && go vet ./...` to verify
4. Submit a PR

## Release Process

When releasing a new version:

1. **Update CHANGELOG.md** on `main`:
   - Rename `## [Unreleased]` to `## [VERSION] - DATE` (e.g., `## [0.2.0] - 2026-04-18`)
   - Add the release date
   - Commit and push to `main`

2. **Create the release**:
   - Go to [Actions → Release → Run workflow](https://github.com/boldandbrad/spin/actions/workflows/release.yml)
   - Enter the version (e.g., `v0.2.0`)
   - Click "Run workflow"

3. **Add release notes**:
   - Copy the changelog section for the new version
   - Paste into the GitHub release body

## Development

See [README.md](./README.md) for usage.