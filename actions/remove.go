package actions

import (
	"fmt"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/lib/efmt"
	"go.polydawn.net/reppl/model"
)

func Remove(c *cli.Context) error {
	tag := c.Args().Get(0)

	p := model.FromFile(".reppl")
	p.DeleteTag(tag)
	p.WriteFile(".reppl")
	fmt.Printf(
		"%s%s %s\n",
		efmt.AnsiWrap("reppl rm", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap(": removed tag", efmt.Ansi_textYellow),
		efmt.AnsiWrap(tag, efmt.Ansi_textYellow, efmt.Ansi_underline),
	)
	return nil
}
