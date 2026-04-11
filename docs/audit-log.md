# Audit Log

Chronological record of audits, releases, documentation passes, and other
maintenance activities. Append-only — newest entries at the bottom.

## 2026-04-11 — /release v0.336.0

- **Commit**: `pending`
- **Outcome**: Released v0.336.0 (darwin-arm64, linux-amd64, linux-arm64). Ships the lodash code-injection CVE fix (GHSA-r5fr-rjxr-66jc) for the docs site dependency tree.
- **Deferred**:
  - 🎯T5 docs/ npm audit is clean — 22 transitive advisories remain in Docusaurus 3.9.2's dep tree (serialize-javascript, picomatch, brace-expansion, path-to-regexp). Tracked as a target; addressed separately via surgical overrides.
