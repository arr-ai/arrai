package rel

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashIdentityDifferential(t *testing.T) {
	if hashIdentity {
		t.Skip("this process is already the hashidentity build")
	}
	t.Run("packages", func(t *testing.T) {
		cmd := exec.Command("go", "test", "-tags", "hashidentity", "-count=1",
			"-skip", "TestNet", "./rel", "./syntax")
		cmd.Dir = ".."
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
	})
	t.Run("corpus", func(t *testing.T) {
		cmd := exec.Command("go", "run", "-tags", "hashidentity", "./cmd/arrai", "test")
		cmd.Dir = ".."
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
	})
}
