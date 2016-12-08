package actions

import (
	"os"

	"github.com/ugorji/go/codec"
	"github.com/urfave/cli"

	"go.polydawn.net/reppl/model"
)

func Show(c *cli.Context) error {
	p := model.FromFile(".reppl")
	enc := codec.NewEncoder(os.Stdout, &codec.JsonHandle{Indent: -1})
	enc.Encode(&p)
	return nil
}
