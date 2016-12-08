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
		fmt.Println("formula already up to date!")
		return nil
	}

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
		if output.Tag != "" {
			p.PutResult(output.Tag, outputName, &rr)
		}
	}

	p.WriteFile(".reppl")

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
				input.Type = ware.Type
			}
		}
	}
	return frm
}

func invokeRepeatr(formulaFileName string) rdef.RunRecord {
	rrBuf := &bytes.Buffer{}
	cmd := Gosh("repeatr", "run", "--ignore-job-exit", formulaFileName,
		Opts{
			Out: rrBuf,
			Err: os.Stderr,
		},
	).Bake()
	cmd.Run()

	fmt.Println(rrBuf.String())
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
