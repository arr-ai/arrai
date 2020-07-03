package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestTransformProtobufToTupleFlow(t *testing.T) {
	code := "//encoding.proto.decode(//os.file('../translate/pb/test/sysl.pb'), 'Module', //os.file('../translate/pb/test/petshop.pb'))"
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(code))
	assert.NoError(t, err)

	codeExpr := pc.CompileExpr(ast)
	val, err := codeExpr.Eval(rel.EmptyScope)

	assert.NoError(t, err)
	tuple, _ := val.(rel.Tuple)
	apps, _ := tuple.Get("apps")

	assert.NotNil(t, apps)
}

func TestTransformProtobufToTupleCompareResult(t *testing.T) {
	code := "//encoding.proto.decode(//os.file('../translate/pb/test/sysl.pb'), 'Module', //os.file('../translate/pb/test/petshop.pb'))"
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(code))
	assert.NoError(t, err)

	codeExpr := pc.CompileExpr(ast)
	val, err := codeExpr.Eval(rel.EmptyScope)

	assert.NoError(t, err)
	tuple, _ := val.(rel.Tuple)

	assert.NotNil(t, tuple)
}
