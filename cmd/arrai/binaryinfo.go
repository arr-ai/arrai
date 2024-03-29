package main

// Binary info variables are dynamically injected via the `-ldflags` flag with `go build`
// Version   - Binary version
// GitFullCommit - Commit SHA of the source code
// BuildDate - Binary build date
// GoVersion - Binary build GoVersion
// BuildOS   - Operating System used to build binary
//
//nolint:gochecknoglobals
var (
	Version       = "unspecified"
	GitFullCommit = "unspecified"
	GitTags       = "unspecified"
	BuildDate     = "unspecified"
	GoVersion     = "unspecified"
	BuildOS       = "unspecified"
	BuildArch     = "unspecified"
)
