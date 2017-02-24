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
	warehouseStrArg_isSet := c.IsSet("warehouse")
	warehouseStr := c.String("warehouse")
	kindStrArg_isSet := c.IsSet("kind")
	kindStr := c.String("kind")

	p := model.FromFile(".reppl")

	wareType := "tar"
	if kindStrArg_isSet {
		wareType = kindStr
	}
	ware := rdef.Ware{
		Type: wareType,
		Hash: hash,
	}
	p.PutManualTag(tag, ware)

	if warehouseStrArg_isSet {
		p.AppendWarehouseForWare(
			ware,
			rdef.WarehouseCoords{rdef.WarehouseCoord(warehouseStr)},
		)
	}

	p.WriteFile(".reppl")

	fmt.Printf(
		"%s %s %s %s %s%s%s\n",
		efmt.AnsiWrap("reppl put", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap("hash:", efmt.Ansi_textYellow),
		efmt.AnsiWrap(tag, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap("=", efmt.Ansi_textYellow),
		efmt.AnsiWrap(ware.Type, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap(":", efmt.Ansi_textYellow),
		efmt.AnsiWrap(ware.Hash, efmt.Ansi_textYellow, efmt.Ansi_underline),
	)
	return nil
}

func PutFile(c *cli.Context) error {
	fmt.Println("put file")
	return nil
}
