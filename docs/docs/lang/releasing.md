---
id: releasing
title: Releasing
---

When a version tag is pushed, it triggers the [Release
workflow](https://github.com/arr-ai/arrai/blob/master/.github/workflows/release.yml)
to publish release binaries to [Arr.ai's GitHub releases
page](https://github.com/arr-ai/arrai/releases).

The release process is automated mostly via
[GoReleaser](https://goreleaser.com/), which creates and deploys
`arrai-X.Y.Z-Os-Arch.tar.gz` and `arrai-X.Y.Z-Windows-Arch.zip` to the [Sysl
Github Release page](https://github.com/arr-ai/arrai/releases). See [GoReleaser
config
file](https://github.com/arr-ai/arrai/blob/master/.github/workflows/.goreleaser.yml)
for further details.

Arr.ai follows a simplified semver model for versioning releases, only ever
incrementing the minor version. The only exception will be a one-time bump to
v1.0.0 when the language is stable enough to offer a backwards compatibility
guarantee.
