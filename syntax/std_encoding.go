package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate/pb"
)

const (
	decodeAttr  = "decode"
	decoderAttr = "decoder"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding",
		stdEncodingCSV(),
		stdEncodingJSON(),
		stdEncodingProto(),
		stdEncodingXlsx(),
		stdEncodingYAML(),
		rel.NewAttr("bytes", mustParseLit(`(decode: \b b)`)),
	)
}

func stdEncodingProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr(decodeAttr, pb.StdProtobufDecoder),
		rel.NewAttr("descriptor", pb.StdProtobufDescriptor),
	)
}

func toDecoderTuple(ctx context.Context, e rel.Expr) (rel.Tuple, error) {
	v, err := e.Eval(ctx, rel.EmptyScope)
	if err != nil {
		return nil, fmt.Errorf("fail to evalute decoder: %v", err)
	}
	if t, isTuple := v.(rel.Tuple); isTuple && t.HasName(decodeAttr) {
		return t, err
	}
	return nil, fmt.Errorf("does not evaluate to a decoder tuple: %v", v)
}

func decode(ctx context.Context, decoder rel.Tuple, bytes []byte) (rel.Value, error) {
	d, err := toDecoderTuple(ctx, decoder)
	if err != nil {
		return nil, err
	}
	decoderFn := d.MustGet(decodeAttr)
	f, isFunction := decoderFn.(rel.Set)
	if !isFunction {
		return nil, fmt.Errorf("does not evaluate to a decoder function: %v", decoderFn)
	}
	return rel.SetCall(ctx, f, rel.NewBytes(bytes))
}
