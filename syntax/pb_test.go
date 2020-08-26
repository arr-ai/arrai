package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const decodePetshop = `let descriptor = //encoding.proto.descriptor(//os.file('../translate/pb/test/sysl.pb'));
	let decodeSyslPb = //encoding.proto.decode(descriptor);
	let shop = decodeSyslPb('Module', //os.file("../translate/pb/test/petshop.pb"));`

func TestTransformProtobufToTupleFlow(t *testing.T) {
	t.Parallel()

	ctx := arraictx.InitRunCtx(context.Background())

	code := decodePetshop + `shop`
	pc := ParseContext{SourceDir: ".."}
	ast, err := pc.Parse(ctx, parser.NewScanner(code))
	assert.NoError(t, err)

	codeExpr, err := pc.CompileExpr(ctx, ast)
	require.NoError(t, err)
	val, err := codeExpr.Eval(ctx, rel.EmptyScope)

	assert.NoError(t, err)
	shop, _ := val.(rel.Tuple)
	apps, _ := shop.Get("apps")

	assert.NotNil(t, apps)
}

func TestTransformProtobufToTupleMapList(t *testing.T) {
	t.Parallel()

	code := decodePetshop + `shop.apps('PetShopApi').name.part.@`
	AssertCodesEvalToSameValue(t, `0`, code)

	code = decodePetshop + `shop.apps('PetShopApi').name.part.@item`
	AssertCodesEvalToSameValue(t, `'PetShopApi'`, code)

	code = decodePetshop + `shop.apps('PetShopApi').attrs('package').s`
	AssertCodesEvalToSameValue(t, "'io.sysl.demo.petshop.api'", code)
	// With index
	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt(0).s`
	AssertCodesEvalToSameValue(t, "'rest'", code)

	code = decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').attrs('patterns').a.elt.@item.s`
	AssertCodesEvalToSameValue(t, "'rest'", code)
}

func TestTransformProtobufToTupleEnum(t *testing.T) {
	t.Parallel()

	code := decodePetshop + `shop.apps('PetShopApi').endpoints('GET /petshop').rest_params.method`
	AssertCodesEvalToSameValue(t, "'GET'", code)

	code = decodePetshop + `shop.apps('PetShopApi').types('Breed').tuple.attr_defs('avgLifespan').primitive`
	AssertCodesEvalToSameValue(t, "'DECIMAL'", code)
}
