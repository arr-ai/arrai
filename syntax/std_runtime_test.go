package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildInfo(t *testing.T) {
	str, err := rel.PrettifyString(GetBuildInfo(
		"DIRTY-v0.99.0", "2020-07-23T10:40:08Z", "d399e13f3670c6698ba35148e6f545322e20e1fb",
		"v0.99.0", "darwin", "amd64", "go1.14 darwin/amd64"), 0)
	assert.Nil(t, err)
	assert.Equal(t, `(
  date: '2020-07-23T10:40:08Z',
  git: (
    commit: 'd399e13f3670c6698ba35148e6f545322e20e1fb',
    tags: ['v0.99.0']
  ),
  go: (
    arch: 'amd64',
    compiler: (
      arch: 'amd64',
      os: 'darwin',
      version: 'go1.14'
    ),
    os: 'darwin'
  ),
  version: 'DIRTY-v0.99.0'
)`, str)
}
