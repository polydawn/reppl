package actions

import (
	"fmt"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/model"
)

func Remove(c *cli.Context) error {
	tag := c.Args().Get(0)

	p := model.FromFile(".reppl")
	p.DeleteTag(tag)
	p.WriteFile(".reppl")
	fmt.Println("removed", tag)
	return nil
}
