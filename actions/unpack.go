package actions

import (
	"fmt"
	"os"

	. "github.com/polydawn/gosh"
	"github.com/urfave/cli"
	rdef "go.polydawn.net/repeatr/api/def"

	"go.polydawn.net/reppl/lib/efmt"
	"go.polydawn.net/reppl/model"
)

func Unpack(c *cli.Context) error {
	// Fathom args.
	tag := c.Args().Get(0)
	unpackPath := c.Args().Get(1)
	warehouseArgStr_isSet := c.IsSet("warehouse")
	warehouseArgStr := c.String("warehouse") // Optional and usually unneeded because our model may also cache warehouses.

	// Load model of project.
	// It should know about the tag and resolve us a ware hash.
	proj := model.FromFile(".reppl")
	ware, err := proj.GetWareByTag(tag)
	if err != nil {
		panic(err)
	}

	// Figure out what warehouseCoords to use, and munge to correct struct.
	// If the was a command line given, prefer that;
	// otherwise, check for one in the project model;
	// if none, we cannot proceed.
	// TODO verify this, actually.  Maybe your cache has it and we can live on that.
	var warehouseCoords rdef.WarehouseCoords
	if warehouseArgStr_isSet {
		warehouseCoords = rdef.WarehouseCoords{rdef.WarehouseCoord(warehouseArgStr)}
	} else {
		warehouseCoords, err = proj.GetWarehousesByWare(ware)
		if err != nil {
			panic(err)
		}
	}

	// Fire all lasers.
	invokeRepeatrUnpack(
		ware,
		unpackPath,
		warehouseCoords,
		false,
	)

	// Confess to our sins.
	fmt.Printf(
		"%s %s %s %s %s %s %s\n",
		efmt.AnsiWrap("reppl unpack", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap("tag:", efmt.Ansi_textYellow),
		efmt.AnsiWrap(tag, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap("=", efmt.Ansi_textYellow),
		efmt.AnsiWrap(ware.Hash, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap("to", efmt.Ansi_textYellow),
		efmt.AnsiWrap(unpackPath, efmt.Ansi_textYellow, efmt.Ansi_underline),
	)
	return nil
}

func invokeRepeatrUnpack(ware rdef.Ware, placePath string, warehouseCoords rdef.WarehouseCoords, skipExists bool) {
	outstrm := efmt.LinePrefixingWriter(
		os.Stderr,
		efmt.AnsiWrap("│ reppl eval >\t", efmt.Ansi_textBrightPurple),
	)
	cmd := Gosh("repeatr", "unpack",
		"--hash", ware.Hash,
		"--kind", ware.Type,
		"--place", placePath,
		"--where", string(warehouseCoords[0]),
		Opts{
			Out: outstrm,
			Err: outstrm,
		},
	).Bake()
	if skipExists {
		cmd = cmd.Bake("--skip-exists")
	}
	cmd.Run()
}
