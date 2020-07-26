package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

// BuildInfo represents arr.ai build information.
var BuildInfo rel.Value

const buildInfoTemplate = `
(
    version: '%s',
    date: '%s',
    git: (
        commit: '%s',
        tags: ['%s'],
    ),
    go: (
        os: '%s',
        arch: '%s',
        compiler: (
            version: '%s',
            os: '%s',
            arch: '%s',
        ),
    ),
)
`

func stdRuntime() rel.Attr {
	return rel.NewTupleAttr("arrai",
		rel.NewAttr("info", BuildInfo),
	)
}

// GetBuildInfo returns arr.ai build information.
func GetBuildInfo(version, date, fullCommit, tags, os, arch, goVersion string) rel.Value {
	expr, err := EvaluateExpr(".",
		fmt.Sprintf(buildInfoTemplate, version, date, fullCommit, tags, os, arch, goVersion, os, arch))
	if err != nil {
		return rel.NewString([]rune(err.Error()))
	}
	return expr
}
