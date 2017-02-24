package actions

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"

	. "github.com/polydawn/gosh"
	"github.com/ugorji/go/codec"
	"github.com/urfave/cli"
	rdef "go.polydawn.net/repeatr/api/def"
	rhitch "go.polydawn.net/repeatr/api/hitch"

	"go.polydawn.net/reppl/lib/efmt"
	"go.polydawn.net/reppl/model"
)

func Eval(c *cli.Context) error {
	// open the formula file
	formulaFileName := c.Args().Get(0)
	pinFileName := formulaFileName + ".pin"
	f, err := os.Open(formulaFileName)
	if err != nil {
		panic("could not open formula")
	}
	defer f.Close()

	// decode the formula file into a formula
	var frm rdef.Formula
	rhitch.DecodeYaml(f, &frm)

	// get our project definition
	p := model.FromFile(".reppl")

	// create the pinned formulas
	pinnedFrm := createPinnedFormula(&p, frm)

	// check if this formula is up to date
	formulaHash := getHash(pinnedFrm)
	if _, exists := p.Memos[formulaHash]; exists {
		fmt.Printf(
			"%s %s %s%s\n",
			efmt.AnsiWrap("┌─", efmt.Ansi_textYellow),
			efmt.AnsiWrap("reppl eval", efmt.Ansi_textBrightYellow),
			efmt.AnsiWrap(formulaFileName, efmt.Ansi_textYellow, efmt.Ansi_underline),
			efmt.AnsiWrap(": no op!  results are on record.", efmt.Ansi_textYellow),
		)
		return nil
	}
	fmt.Printf(
		"%s %s %s%s\n",
		efmt.AnsiWrap("┌─", efmt.Ansi_textYellow),
		efmt.AnsiWrap("reppl eval", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap(formulaFileName, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap(": looks new, no memoized result!  evaluating...", efmt.Ansi_textYellow),
	)

	// write the pinned formula file as JSON
	writeFormula(&pinnedFrm, pinFileName)

	// make repeatr go now!
	rr := invokeRepeatr(pinFileName)

	// add the formula hash to the run record
	rr.FormulaHID = formulaHash
	// add the run record hash to the run record
	rr.HID = getHash(rr)

	// save tagged outputs
	for outputName, output := range frm.Outputs {
		if output.Tag == "" {
			continue
		}
		p.PutResult(output.Tag, outputName, &rr)
	}

	// memorize all the warehouses that were listed as destinations for outputs
	for outputName, output := range frm.Outputs {
		if output.Tag == "" {
			continue
		}
		p.AppendWarehouseForWare(rr.Results[outputName].Ware, output.Warehouses)
	}

	p.WriteFile(".reppl")
	fmt.Printf(
		"%s %s %s%s\n",
		efmt.AnsiWrap("└─", efmt.Ansi_textYellow),
		efmt.AnsiWrap("reppl eval", efmt.Ansi_textBrightYellow),
		efmt.AnsiWrap(formulaFileName, efmt.Ansi_textYellow, efmt.Ansi_underline),
		efmt.AnsiWrap(": done!  results saved.", efmt.Ansi_textYellow),
	)

	return nil
}

func getHash(v interface{}) string {
	hash := sha512.New384()
	enc := codec.NewEncoder(hash, &codec.JsonHandle{})
	err := enc.Encode(v)
	if err != nil {
		panic("could not hash struct")
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func createPinnedFormula(p *model.Project, frm rdef.Formula) rdef.Formula {
	// add our hashes by tags
	for _, input := range frm.Inputs {
		if input.Tag != "" {
			ware, err := p.GetWareByTag(input.Tag)
			if err == nil {
				input.Hash = ware.Hash
				if input.Type == "" {
					input.Type = ware.Type
				}
			}
		}
	}
	// append any warehouses we know of
	for _, input := range frm.Inputs {
		ware := rdef.Ware{input.Type, input.Hash}
		moreWarehouseCoords, err := p.GetWarehousesByWare(ware)
		if err != nil {
			// nbd if we don't have any.  hope the formula had some of its own; but if not, that error isn't for our layer to raise.
			continue
		}
		input.Warehouses = append(input.Warehouses, moreWarehouseCoords...)
	}
	return frm
}

func invokeRepeatr(formulaFileName string) rdef.RunRecord {
	rrBuf := &bytes.Buffer{}
	cmd := Gosh("repeatr", "run", "--ignore-job-exit", formulaFileName,
		Opts{
			Out: rrBuf,
			Err: efmt.LinePrefixingWriter(
				os.Stderr,
				efmt.AnsiWrap("│ reppl eval >\t", efmt.Ansi_textBrightPurple),
			),
		},
	).Bake()
	cmd.Run()

	fmt.Fprintln(efmt.LinePrefixingWriter(
		os.Stderr,
		efmt.AnsiWrap("│ reppl eval ∴⟩\t", efmt.Ansi_textYellow),
	), rrBuf.String())
	var rr rdef.RunRecord
	dec := codec.NewDecoder(rrBuf, &codec.JsonHandle{})
	err := dec.Decode(&rr)
	if err != nil {
		panic("error reading run record: " + err.Error())
	}
	return rr
}

func writeFormula(frm *rdef.Formula, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic("error opening pin file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	enc := codec.NewEncoder(w, &codec.JsonHandle{Indent: -1})
	err = enc.Encode(frm)
	if err != nil {
		panic("could not write pin file")
	}
	w.Write([]byte{'\n'})
}
