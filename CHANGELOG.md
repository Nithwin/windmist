# Changelog

All notable changes to the WindMist CLI project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- **Repository Infrastructure**: Standardized GitHub issue forms (`bug_report.yml`, `feature_request.yml`, `documentation.yml`, `config.yml`) and `pull_request_template.md`.
- **Labels Infrastructure**: Automated label specification (`labels.json`) and synchronization script (`sync_labels.sh`).
- **Automated CI**: Fast, single-purpose PR validation workflow ([`ci.yml`](.github/workflows/ci.yml)) executing `gofmt`, `go vet`, `go test -race -timeout 10m`, and binary compilation checks on Ubuntu Linux.
- **Release Automation**: Multi-platform cross-compilation (`linux`, `macOS`, `windows` across `amd64` and `arm64`), `tar.gz`/`zip` + standalone `.exe` packaging, and SHA256 checksum generation (`checksums.txt`) using GoReleaser ([`.goreleaser.yaml`](.goreleaser.yaml) and [`release.yml`](.github/workflows/release.yml)).
- **Dynamic Build Metadata**: Enabled `-X` linker flag injection for `Version`, `Commit`, and `Date` in `cmd/version.go`.
- **Self-Uninstallation Command**: Added built-in `windmist uninstall` (`cmd/uninstall.go`) with flags `-y/--yes` and `-p/--purge` to allow the CLI binary to self-remove from disk cleanly.
- **Universal Installer Script**: Created `scripts/install.sh` for 1-line `curl | bash` installation across Linux and macOS.

### Changed
- Normalized formatting across all 38 Go packages (`cmd/` and `internal/`) to adhere to `gofmt`.
- Updated `go.mod` language specification to `go 1.26`.

---

## [v1.0.0] - 2026-07-14

### Added
- Initial release of the WindMist CLI software engineering agent with modular AI provider support (Gemini), filesystem editing tools, and Lip Gloss/Glamour TUI rendering.
