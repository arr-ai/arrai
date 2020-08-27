package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/pkg/errors"
	"go/parser"
	"go/token"
)

func stdLang() rel.Attr {
	return rel.NewTupleAttr(
		"lang", rel.NewTupleAttr(
			"go", rel.NewNativeFunctionAttr(
				"parse", func(v rel.Value) (rel.Value, error) {
					// Accepts []byte because ParseFile ultimately uses []byte.
					doParse := func(src []byte) (rel.Value, error) {
						fset := token.NewFileSet()
						f, err := parser.ParseFile(fset, "", src, 0)
						if err != nil {
							return nil, err
						}
						return rel.NewValue(f)
					}

					switch t := v.(type) {
					case rel.GenericSet:
						if !t.IsTrue() {
							return rel.None, nil
						}
					case rel.Bytes:
						return doParse(t.Bytes())
					case rel.String:
						return doParse([]byte(t.String()))
					}
					return nil, errors.Errorf("parse requires string or bytes, not %T", v)
				})))
}
