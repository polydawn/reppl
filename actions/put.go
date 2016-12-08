package actions

import (
	"fmt"

	"github.com/urfave/cli"
	rdef "go.polydawn.net/repeatr/api/def"

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
	fmt.Println("put", tag, "=", hash)
	return nil
}

func PutFile(c *cli.Context) error {
	fmt.Println("put file")
	return nil
}
