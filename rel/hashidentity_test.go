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
	cmd := exec.Command("go", "test", "-tags", "hashidentity", "-count=1",
		"-run", "TestAlgebraEqualImpliesHash128|TestAlgebraDictVsGenericSet", ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
}
