package buildinfo

import (
	"context"
	"time"
)

// BuildData is a struct that contains build information of arr.ai
type BuildData struct {
	Version, Date, FullCommit, Tags, Os, Arch, GoVersion string
}

type buildDataKey int

const buildInfoKey buildDataKey = iota

// BuildInfo represents arr.ai build information.
var BuildInfo BuildData

// WithBuildData puts build data into context.
func WithBuildData(ctx context.Context, b BuildData) context.Context {
	return context.WithValue(ctx, buildInfoKey, b)
}

// WithPackageBuildData inserts the global variable BuildInfo into context. BuildInfo is initialized by the main
// package.
func WithPackageBuildData(ctx context.Context) context.Context {
	return WithBuildData(ctx, BuildInfo)
}

// BuildDataFrom takes in context and return the BuildData.
func BuildDataFrom(ctx context.Context) *BuildData {
	if b := ctx.Value(buildInfoKey); b != nil {
		if data, is := b.(BuildData); is {
			return &data
		}
	}
	return nil
}

// BuildDateFrom takes in context and return the associated build date. It returns nil if build data is not in context.
func BuildDateFrom(ctx context.Context) (_ *time.Time, err error) {
	var buildDate time.Time
	buildData := BuildDataFrom(ctx)

	switch {
	case buildData == nil || buildData.Date == "" || buildData.Date == "unspecified":
		// if buildVersion is not found, just provide warning since the Deprecate function is still called.
		return nil, nil
	default:
		buildDate, err = time.Parse(time.RFC3339, buildData.Date)
		if err != nil {
			// unlikely to happen since build date is automatically generated
			return nil, err
		}
	}
	return &buildDate, nil
}

// SetBuildInfo sets the global BuildInfo variable in the buildinfo package.
func SetBuildInfo(version, date, fullCommit, tags, os, arch, goVersion string) {
	BuildInfo = BuildData{
		Version:    version,
		Date:       date,
		FullCommit: fullCommit,
		Tags:       tags,
		Os:         os,
		Arch:       arch,
		GoVersion:  goVersion,
	}
}
