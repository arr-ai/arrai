package mod

import (
	"github.com/joshcarp/gop/gop/cli"
	"github.com/spf13/afero"
	"log"
	"os"
)

func New(fs afero.Fs)cli.Retriever{
	tokenmap, _ := cli.NewTokenMap("ARRAI_TOKENS", "GIT_CREDENTIALS")
	return cli.Moduler(fs, "arrai_modules.yaml","arrai_modules", os.Getenv("ARRAI_PROXY"), tokenmap, log.Printf)
}
