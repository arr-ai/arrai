package syntax

import (
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// BuildInfo represents arr.ai build information.
var BuildInfo rel.Value

func stdRuntime() rel.Attr {
	return rel.NewTupleAttr("arrai",
		rel.NewAttr("info", BuildInfo),
	)
}

// GetBuildInfo returns arr.ai build information.
func GetBuildInfo(version, date, fullCommit, tags, os, arch, goVersion string) rel.Value {
	gitInfo := rel.NewTuple(
		rel.NewAttr("commit", rel.NewString([]rune(fullCommit))),
		// param tags has only one tag now.
		rel.NewAttr("tags", rel.NewArray(rel.NewString([]rune(tags)))))

	// Extract 'go1.14' from 'go1.14 darwin/amd64'
	infos := strings.Split(goVersion, " ")
	var goVersionNo rel.Value
	if len(infos) > 0 {
		goVersionNo = rel.NewString([]rune(infos[0]))
	} else {
		goVersionNo = rel.NewString([]rune(goVersion))
	}

	compiler := rel.NewTuple(
		rel.NewAttr("version", goVersionNo),
		rel.NewAttr("os", rel.NewString([]rune(os))),
		rel.NewAttr("arch", rel.NewString([]rune(arch))))

	goInfo := rel.NewTuple(
		rel.NewAttr("os", rel.NewString([]rune(os))),
		rel.NewAttr("arch", rel.NewString([]rune(arch))),
		rel.NewAttr("compiler", compiler))

	info := rel.NewTuple(
		rel.NewAttr("version", rel.NewString([]rune(version))),
		rel.NewAttr("date", rel.NewString([]rune(date))),
		rel.NewAttr("git", gitInfo),
		rel.NewAttr("go", goInfo),
	)

	return rel.NewBuildInfoTuple(info)
}
