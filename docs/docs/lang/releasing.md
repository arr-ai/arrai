---
id: releasing
title: Releasing
---

A release is cut by pushing a version tag (`vX.Y.0`) to `master`. Tags are
created deliberately by a maintainer, never automatically on merge.

Pushing the tag triggers the [Release
workflow](https://github.com/arr-ai/arrai/blob/master/.github/workflows/release.yml),
which runs [GoReleaser](https://goreleaser.com/) to build
`arrai_vX.Y.0_<os>-<arch>.tar.gz` (and `.zip` for Windows) and publish them,
with checksums, to [Arr.ai's GitHub releases
page](https://github.com/arr-ai/arrai/releases). See the [GoReleaser config
file](https://github.com/arr-ai/arrai/blob/master/.github/goreleaser_configs/goreleaser.yml)
for details.

Arr.ai follows a simplified semver model for versioning releases, only ever
incrementing the minor version. The only exception will be a one-time bump to
v1.0.0 when the language is stable enough to offer a backwards compatibility
guarantee.
