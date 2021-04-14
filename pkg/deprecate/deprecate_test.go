package deprecate

import (
	"context"
	"fmt"
	"testing"

	"github.com/arr-ai/arrai/pkg/buildinfo"
	"github.com/arr-ai/wbnf/parser"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewDeprecator(t *testing.T) {
	t.Parallel()
	_, err := NewDeprecator("", "2000-01-01", "2000-04-01", "2000-02-01")
	assert.EqualError(t,
		err, `weak, strong, and crash versions are not sorted: "2000-01-01", "2000-04-01", "2000-02-01"`,
	)
}

func TestDeprecate(t *testing.T) {
	// FIXME: use hook attached to a local logger instead of the global one
	hook := test.NewGlobal()
	deprecator := MustNewDeprecator("feature", "2000-01-02", "2000-01-04", "2000-01-06")
	scanner := parser.NewScanner("expression")
	scannerStr := scanner.Context(parser.DefaultLimit)

	errMsg := func(msg string) string {
		return fmt.Sprintf("%s\n%s", msg, scannerStr)
	}

	assert.NoError(t, deprecator.Deprecate(context.Background(), *scanner))
	assert.Equal(t, errMsg("feature is being deprecated"), hook.LastEntry().Message)
	hook.Reset()

	ctx := buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: "unspecified"})
	assert.NoError(t, deprecator.Deprecate(ctx, *scanner))
	assert.Equal(t, errMsg("feature is being deprecated"), hook.LastEntry().Message)
	hook.Reset()

	ctx = buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: ""})
	assert.NoError(t, deprecator.Deprecate(ctx, *scanner))
	assert.Equal(t, errMsg("feature is being deprecated"), hook.LastEntry().Message)
	hook.Reset()

	ctx = buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: "2000-01-01T00:00:00Z"})
	assert.NoError(t, deprecator.Deprecate(ctx, *scanner))
	assert.Equal(t, 0, len(hook.Entries))
	hook.Reset()

	ctx = buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: "2000-01-03T00:00:00Z"})
	assert.NoError(t, deprecator.Deprecate(ctx, *scanner))
	assert.Equal(t, errMsg("feature is being deprecated"), hook.LastEntry().Message)
	hook.Reset()

	ctx = buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: "2000-01-05T00:00:00Z"})
	assert.NoError(t, deprecator.Deprecate(ctx, *scanner))
	assert.Equal(t,
		errMsg(fmt.Sprintf("feature is being deprecated (pausing %s...)", delayDurationStr)),
		hook.LastEntry().Message,
	)
	hook.Reset()

	ctx = buildinfo.WithBuildData(context.Background(), buildinfo.BuildData{Date: "2000-01-07T00:00:00Z"})
	assert.EqualError(t, deprecator.Deprecate(ctx, *scanner), "feature is deprecated")
	assert.Equal(t, 0, len(hook.Entries))
	hook.Reset()
}
