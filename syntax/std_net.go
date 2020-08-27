package syntax

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

func stdNet() rel.Attr {
	return rel.NewTupleAttr(
		"net",
		rel.NewTupleAttr(
			"http",
			rel.NewNativeFunctionAttr("get", func(_ context.Context, v rel.Value) (rel.Value, error) {
				url, is := tools.ValueAsString(v)
				if !is {
					return nil, fmt.Errorf("//net.http.get: url not a string: %v", url)
				}
				resp, err := http.Get(url) //nolint:gosec
				if err != nil {
					return nil, err
				}

				buf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, err
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
				), nil
			}),
		),
	)
}
