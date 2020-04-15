package syntax

import (
	"io/ioutil"
	"net/http"

	"github.com/arr-ai/arrai/rel"
)

func stdNet() rel.Attr {
	return rel.NewTupleAttr(
		"net",
		rel.NewTupleAttr(
			"http",
			rel.NewNativeFunctionAttr("get", func(v rel.Value) rel.Value {
				resp, err := http.Get(mustAsString(v))
				if err != nil {
					panic(err)
				}

				buf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				header := rel.NewTuple()
				for k, vs := range resp.Header {
					vals := make([]rel.Value, 0, len(vs))
					for _, v := range vs {
						vals = append(vals, rel.NewString([]rune(v)))
					}
					header = header.With(k, rel.NewArray(vals...))
				}

				return rel.NewTuple(
					rel.NewAttr("status", rel.NewString([]rune(resp.Status))),
					rel.NewAttr("status_code", rel.NewNumber(float64(resp.StatusCode))),
					rel.NewAttr("header", header),
					rel.NewAttr("body", rel.NewBytes(buf)),
				)
			}),
		),
	)
}
