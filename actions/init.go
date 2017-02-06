package actions

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/lib/efmt"
	"go.polydawn.net/reppl/model"
)

func Init(c *cli.Context) error {
	p := model.Project{}

	if _, err := os.Stat(".reppl"); os.IsNotExist(err) {
		p.Init()
		p.WriteFile(".reppl")
		fmt.Printf(
			"%s%s\n",
			efmt.AnsiWrap("reppl init", efmt.Ansi_textBrightYellow),
			efmt.AnsiWrap(": created new project file!", efmt.Ansi_textYellow),
		)
		return nil
	} else {
		fmt.Printf(
			"%s%s\n",
			efmt.AnsiWrap("reppl init", efmt.Ansi_textBrightYellow),
			efmt.AnsiWrap(": project file aready exists.", efmt.Ansi_textYellow),
		)
		return nil
	}
}
