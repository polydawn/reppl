package actions

import (
	"fmt"

	"github.com/urfave/cli"
	rdef "go.polydawn.net/repeatr/api/def"

	"go.polydawn.net/reppl/lib/efmt"
	"go.polydawn.net/reppl/model"
)

func PutHash(c *cli.Context) error {
	tag := c.Args().Get(0)
	hash := c.Args().Get(1)

	p := model.FromFile(".reppl")
	ware := rdef.Ware{
		Type: "tar",
		Hash: hash,
	}
	p.PutManualTag(tag, ware)
	p.WriteFile(".reppl")
	fmt.Printf(
		"%s %s %s %s %s\n",
		efmt.AnsiWrap("reppl put", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap("hash", efmt.Ansi_textYellow),
		efmt.AnsiWrap(tag, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap("=", efmt.Ansi_textYellow),
		efmt.AnsiWrap(hash, efmt.Ansi_textYellow, efmt.Ansi_underline),
	)
	return nil
}

func PutFile(c *cli.Context) error {
	fmt.Println("put file")
	return nil
}
