package syntax

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

type httpConfig struct {
	header map[string][]string
}

func stdNet() rel.Attr {
	return rel.NewTupleAttr(
		"net",
		rel.NewTupleAttr(
			"http",
			createFunc2Attr("get", func(_ context.Context, configArg, urlArg rel.Value) (rel.Value, error) {
				config, err := parseConfig(configArg)
				if err != nil {
					return nil, err
				}
				url, err := parseURL(urlArg)
				if err != nil {
					return nil, err
				}
				return get(url, config.header)
			}),
			createFunc3Attr("post",
				func(_ context.Context, configArg, urlArg, bodyArg rel.Value) (rel.Value, error) {
					config, err := parseConfig(configArg)
					if err != nil {
						return nil, err
					}
					url, err := parseURL(urlArg)
					if err != nil {
						return nil, err
					}
					body, err := parseBody(bodyArg)
					if err != nil {
						return nil, err
					}
					return post(url, config.header, body)
				}),
		),
	)
}

// send sends a request of type method to url with headers and body and returns a value wrapping the
// response.
func send(method, url string, headers map[string][]string, body io.Reader) (rel.Value, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if len(headers) > 0 {
		req.Header = headers
	}

	res, err := http.DefaultClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return parseResponse(res)
}

// get sends a GET request and returns a value wrapping the response.
func get(url string, headers map[string][]string) (rel.Value, error) {
	return send("GET", url, headers, strings.NewReader(""))
}

// post sends a POST request and returns a value wrapping the response.
func post(url string, headers map[string][]string, body io.Reader) (rel.Value, error) {
	return send("POST", url, headers, body)
}

// parseConfig returns the config arg as a httpConfig.
func parseConfig(configArg rel.Value) (*httpConfig, error) {
	config, ok := configArg.(*rel.GenericTuple)
	if !ok {
		return nil, errors.Errorf("first arg (config) must be tuple, not %s", rel.ValueTypeAsString(configArg))
	}
	head, ok := config.Get("header")
	if !ok {
		return &httpConfig{}, nil
	}

	header, err := parseHeader(head)
	if err != nil {
		return nil, err
	}

	return &httpConfig{header: header}, nil
}

// parseHeader returns the header of the config arg as a map.
func parseHeader(header rel.Value) (map[string][]string, error) {
	headDict, ok := header.(rel.Dict)
	if !ok {
		return nil, errors.Errorf("header must be a dict, not %s", rel.ValueTypeAsString(headDict))
	}

	out := map[string][]string{}
	for e := headDict.DictEnumerator(); e.MoveNext(); {
		kv, vv := e.Current()
		k, ok := tools.ValueAsString(kv)
		if !ok {
			return nil, errors.Errorf("header keys must be strings, not %s", rel.ValueTypeAsString(kv))
		}
		switch t := vv.(type) {
		case rel.String:
			out[k] = []string{t.String()}
		case rel.Array:
			vs := make([]string, t.Count())
			for _, val := range t.Values() {
				valStr, is := tools.ValueAsString(val)
				if !is {
					return nil, errors.Errorf(
						"header values must be strings or string arrays, not arrays of %s", rel.ValueTypeAsString(val))
				}
				vs = append(vs, valStr)
			}
			out[k] = vs
		default:
			return nil, errors.Errorf("header values must be strings or string arrays, not %s", rel.ValueTypeAsString(vv))
		}
	}
	return out, nil
}

// parseURL returns the URL arg as a string.
func parseURL(urlArg rel.Value) (string, error) {
	url, is := tools.ValueAsString(urlArg)
	if !is {
		return "", errors.Errorf("second arg (url) must be a string, not %s", rel.ValueTypeAsString(urlArg))
	}
	return url, nil
}

// parseBody returns the body arg as a Reader.
func parseBody(bodyArg rel.Value) (io.Reader, error) {
	body, is := tools.ValueAsString(bodyArg)
	if is {
		return strings.NewReader(body), nil
	}

	bodyBytes, is := tools.ValueAsBytes(bodyArg)
	if is {
		return bytes.NewReader(bodyBytes), nil
	}

	return nil, errors.Errorf("third arg (body) must be a string or bytes, not %s", rel.ValueTypeAsString(bodyArg))
}

// parseResponse parses an HTTP response into an arr.ai value.
func parseResponse(resp *http.Response) (rel.Value, error) {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	entries := make([]rel.DictEntryTuple, len(resp.Header))
	for k, vs := range resp.Header {
		vals := make([]rel.Value, 0, len(vs))
		for _, v := range vs {
			vals = append(vals, rel.NewString([]rune(v)))
		}
		entries = append(entries, rel.NewDictEntryTuple(rel.NewString([]rune(k)), rel.NewArray(vals...)))
	}
	header := rel.MustNewDict(false, entries...)

	return rel.NewTuple(
		rel.NewAttr("status", rel.NewString([]rune(resp.Status))),
		rel.NewAttr("status_code", rel.NewNumber(float64(resp.StatusCode))),
		rel.NewAttr("header", header),
		rel.NewAttr("body", rel.NewBytes(buf)),
	), nil
}
