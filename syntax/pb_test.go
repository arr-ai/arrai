package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestTransformProtobufToTupleFlow(t *testing.T) {
	code := "//encoding.proto.decode(//os.file('../translate/pb/test/sysl.pb'), 'Module', " +
		"//os.file('../translate/pb/test/petshop.pb'))"
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(parser.NewScanner(code))
	assert.NoError(t, err)

	codeExpr := pc.CompileExpr(ast)
	val, err := codeExpr.Eval(rel.EmptyScope)

	assert.NoError(t, err)
	shop, _ := val.(rel.Tuple)
	apps, _ := shop.Get("apps")

	assert.NotNil(t, apps)
}

func TestTransformProtobufToTupleCompareResult(t *testing.T) {
	code := `let shop = //encoding.proto.decode(//os.file('../translate/pb/test/sysl.pb'), 'Module', //os.file('../translate/pb/test/petshop.pb'));` +
		`shop.apps('PetShopApi').attrs('package').s`
	AssertCodesEvalToSameValue(t, "'io.sysl.demo.petshop.api'", code)
	// code = `let shop = //encoding.proto.decode(//os.file('../translate/pb/test/sysl.pb'), 'Module', //os.file('../translate/pb/test/petshop.pb'));` +
	// 	`shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a`
	// AssertCodesEvalToSameValue(t, `'(elt: [(@: 0, @item: (s: \'rest\'))])'`, code)
}
