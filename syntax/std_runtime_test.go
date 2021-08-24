package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/buildinfo"
	"github.com/arr-ai/arrai/rel"
)

func TestGetBuildInfo(t *testing.T) {
	assert.Equal(t, `(date: '2020-07-23T10:40:08Z', `+
		`git: (commit: 'd399e13f3670c6698ba35148e6f545322e20e1fb', tags: {'v0.99.0'}), `+
		`go: (arch: 'amd64', compiler: (arch: 'amd64', os: 'darwin', version: 'go1.15 darwin/amd64'), os: 'darwin'), `+
		`version: 'DIRTY-v0.99.0')`,
		rel.Repr(
			GetBuildInfo(
				buildinfo.BuildData{
					Version:    "DIRTY-v0.99.0",
					Date:       "2020-07-23T10:40:08Z",
					FullCommit: "d399e13f3670c6698ba35148e6f545322e20e1fb",
					Tags:       "v0.99.0",
					Os:         "darwin",
					Arch:       "amd64",
					GoVersion:  "go1.15 darwin/amd64",
				},
			),
		),
	)
}
