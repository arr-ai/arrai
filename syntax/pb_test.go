package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

const decodePetshop = `let sysl = //encoding.proto.decode(//encoding.proto.proto,
	//os.file('../translate/pb/test/sysl.pb'));
	let decodeSyslPb = //encoding.proto.decode(sysl);
	let shop = decodeSyslPb('Module', //os.file("../translate/pb/test/petshop.pb"));`

func TestTransformProtobufToTupleFlow(t *testing.T) {
	t.Parallel()

	code := decodePetshop + `shop`
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

func TestTransformProtobufToTupleMapList(t *testing.T) {
	t.Parallel()

	code := decodePetshop + `shop.apps('PetShopApi').name.part.@`
	AssertCodesEvalToSameValue(t, `0`, code)

	code = decodePetshop + `shop.apps('PetShopApi').name.part.@item.@item`
	AssertCodesEvalToSameValue(t, `'PetShopApi'`, code)

	code = decodePetshop + `shop.apps('PetShopApi').attrs('package').s`
	AssertCodesEvalToSameValue(t, "'io.sysl.demo.petshop.api'", code)

	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt(0).@item.s`
	AssertCodesEvalToSameValue(t, "'rest'", code)

	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt.@item.@item.s`
	AssertCodesEvalToSameValue(t, "'rest'", code)

	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt(0).@`
	AssertCodesEvalToSameValue(t, "0", code)

	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt.@`
	AssertCodesEvalToSameValue(t, "0", code)
}

func TestTransformProtobufToTupleEnum(t *testing.T) {
	t.Parallel()

	code := decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').rest_params.method`
	AssertCodesEvalToSameValue(t, "'GET'", code)

	code = decodePetshop + `shop.apps('PetShopApi').types('Breed').tuple.attr_defs('avgLifespan').primitive`
	AssertCodesEvalToSameValue(t, "'DECIMAL'", code)
}
