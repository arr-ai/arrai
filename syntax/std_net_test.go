package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetGet(t *testing.T) {
	t.Parallel()

	expectedBody := rel.NewBytes([]byte(`all: test lint wasm

# TODO: If this Makefile is ever used for CI, suppress timingsensitive there.
test:
	go test $(GOTESTFLAGS) -tags timingsensitive ./...

lint:
	golangci-lint run

wasm:
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai
`))
	expectedStatus := rel.NewString([]rune("200 OK"))
	expectedStatusCode := rel.NewNumber(float64(200))
	expectedContentType := rel.NewString([]rune("text/plain; charset=utf-8"))

	result, err := EvaluateExpr(
		"",
		`//.net.http.get("https://raw.githubusercontent.com/arr-ai/arrai/cf1326f7b61178e3e98aff30540e10cb73449445/Makefile")`,
	)
	require.NoError(t, err)

	assert.True(t, expectedBody.Equal(result.(rel.Tuple).MustGet("body")))
	assert.True(t, expectedStatus.Equal(result.(rel.Tuple).MustGet("status")))
	assert.True(t, expectedStatusCode.Equal(result.(rel.Tuple).MustGet("status_code")))
	assert.True(t, rel.
		NewArray(expectedContentType).
		Equal(
			result.(rel.Tuple).
				MustGet("header").(rel.Tuple).
				MustGet("Content-Type"),
		),
	)
}
