## Project knowledge

Thie repository contains a **Grafana plugin**. You must Read @./.config/AGENTS/instructions.md before doing changes.

## Release process

Releases are driven by semantic-release from `.github/workflows/release.yml` and `.releaserc.json`. Use semantic commits / Conventional Commits for changes so release versions and changelog entries are generated correctly, for example `fix: ...`, `feat: ...`, and `feat!: ...` or `BREAKING CHANGE:` for breaking changes.
