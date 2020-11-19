package syntax

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

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
					return nil, errors.Errorf("//net.http.get: url not a string: %v", url)
				}
				return get(url)
			}),
			createFunc2Attr("post", func(_ context.Context, config rel.Value, v rel.Value) (rel.Value, error) {
				c, ok := config.(*rel.GenericTuple)
				if !ok {
					return nil, errors.Errorf("//net.http.post: first arg must be tuple, not %s", rel.ValueTypeAsString(config))
				}
				url, is := tools.ValueAsString(v)
				if !is {
					return nil, errors.Errorf("//net.http.post: url not a string: %v", url)
				}

				var contentType string
				if ct, ok := c.Get("contentType"); !ok {
					contentType = "text/plan"
				} else {
					contentType = ct.String()
				}

				var body string
				if b, ok := c.Get("body"); !ok {
					body = ""
				} else {
					body = b.String()
				}
				return post(url, contentType, body)
			}),
		),
	)
}

// get sends a GET request and returns a value wrapping the response.
func get(url string) (rel.Value, error) {
	r, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return parseResponse(r)
}

// post sends a POST request and returns a value wrapping the response.
func post(url, contentType, body string) (rel.Value, error) {
	r, err := http.Post(url, contentType, strings.NewReader(body)) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return parseResponse(r)
}

// parseResponse parses an HTTP response into an arr.ai value.
func parseResponse(resp *http.Response) (rel.Value, error) {
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
}
