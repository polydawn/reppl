package model

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ugorji/go/codec"
	rdef "go.polydawn.net/repeatr/api/def"
)

type Project struct {
	Tags        map[string]ReleaseRecord           // map tag->{ware,backstory}
	RunRecords  map[string]*rdef.RunRecord         // map rrhid->rr
	Memos       map[string]string                  // index frmhid->rrhid
	Whereabouts map[rdef.Ware]rdef.WarehouseCoords // map ware->warehousecoords
}

type projectOutput struct {
	Tags        map[string]ReleaseRecord
	RunRecords  map[string]*rdef.RunRecord
	Memos       map[string]string
	Whereabouts map[string]rdef.WarehouseCoords
}

type ReleaseRecord struct {
	Ware         rdef.Ware
	RunRecordHID string // blank if a tag was manual
}

func (p *Project) Init() {
	p.Tags = make(map[string]ReleaseRecord)
	p.RunRecords = make(map[string]*rdef.RunRecord)
	p.Memos = make(map[string]string)
	p.Whereabouts = make(map[rdef.Ware]rdef.WarehouseCoords)
}

func (p *Project) WriteFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("error opening project file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	whereabouts := make(map[string]rdef.WarehouseCoords)
	for ware, coords := range p.Whereabouts {
		wareString := fmt.Sprintf("%s:%s", ware.Type, ware.Hash)
		whereabouts[wareString] = coords
	}

	enc := codec.NewEncoder(w, &codec.JsonHandle{Indent: -1})
	err = enc.Encode(&projectOutput{
		Tags:        p.Tags,
		RunRecords:  p.RunRecords,
		Memos:       p.Memos,
		Whereabouts: whereabouts,
	})
	if err != nil {
		panic("could not write project file")
	}
	w.Write([]byte{'\n'})
}

func FromFile(filename string) Project {
	f, err := os.Open(filename)
	if err != nil {
		panic("error opening project file")
	}
	defer f.Close()

	r := bufio.NewReader(f)
	intermediate := projectOutput{}
	dec := codec.NewDecoder(r, &codec.JsonHandle{})
	err = dec.Decode(&intermediate)
	if err != nil {
		panic("error reading project file")
	}

	whereabouts := make(map[rdef.Ware]rdef.WarehouseCoords)
	for wareStr, coords := range intermediate.Whereabouts {
		parts := strings.SplitN(wareStr, ":", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintf("Invalid ware key string: %s\n", parts))
		}
		ware := rdef.Ware{Type: parts[0], Hash: parts[1]}
		whereabouts[ware] = coords
	}
	p := Project{
		Tags:        intermediate.Tags,
		RunRecords:  intermediate.RunRecords,
		Memos:       intermediate.Memos,
		Whereabouts: whereabouts,
	}

	return p
}

func (p *Project) PutManualTag(tag string, ware rdef.Ware) {
	_, hadPrev := p.Tags[tag]
	p.Tags[tag] = ReleaseRecord{ware, ""}
	if hadPrev {
		p.retainFilter()
	}
}

func (p *Project) AppendWarehouseForWare(ware rdef.Ware, moreCoords rdef.WarehouseCoords) {
	coords, _ := p.Whereabouts[ware]
	// Append, putting the most recent ones first.
	coords = append(moreCoords, coords...)
	// Filter out any duplicates.
	set := make(map[rdef.WarehouseCoord]struct{})
	n := 0
	for i, v := range coords {
		_, exists := set[v]
		if exists {
			continue // leave it behind
		}
		set[v] = struct{}{}
		coords[n] = coords[i]
	}
	p.Whereabouts[ware] = coords[0 : n+1]
}

func (p *Project) DeleteTag(tag string) {
	_, hadPrev := p.Tags[tag]
	if hadPrev {
		delete(p.Tags, tag)
		p.retainFilter()
	}
}

func (p *Project) GetWareByTag(tag string) (rdef.Ware, error) {
	_, exists := p.Tags[tag]
	if exists {
		return p.Tags[tag].Ware, nil
	} else {
		return rdef.Ware{}, errors.New("not found")
	}
}

func (p *Project) GetWarehousesByWare(ware rdef.Ware) (rdef.WarehouseCoords, error) {
	coords, exists := p.Whereabouts[ware]
	if exists {
		return coords, nil
	} else {
		return nil, fmt.Errorf("no warehouses known for ware %s:%s", ware.Type, ware.Hash)
	}
}

func (p *Project) PutResult(tag string, resultName string, rr *rdef.RunRecord) {
	p.Tags[tag] = ReleaseRecord{rr.Results[resultName].Ware, rr.HID}
	p.RunRecords[rr.HID] = rr
	p.Memos[rr.FormulaHID] = rr.HID
	p.retainFilter()
}

func (p *Project) retainFilter() {
	// "Sweep".  (The `Tags` map is the marks.)
	oldRunRecords := p.RunRecords
	oldWhereabouts := p.Whereabouts
	p.RunRecords = make(map[string]*rdef.RunRecord)
	p.Memos = make(map[string]string)
	p.Whereabouts = make(map[rdef.Ware]rdef.WarehouseCoords)
	// Rebuild `RunRecords` by whitelisting prev values still ref'd by `Tags`.
	for tag, release := range p.Tags {
		if release.RunRecordHID == "" {
			continue // skip.  it's just a fiat release; doesn't ref anything.
		}
		runRecord, ok := oldRunRecords[release.RunRecordHID]
		if !ok {
			panic(fmt.Errorf("db integrity violation: dangling runrecord -- release %q points to %q", tag, release.RunRecordHID))
		}
		p.RunRecords[release.RunRecordHID] = runRecord
	}
	// Rebuild `Memos` index from `RunRecords`.
	for _, runRecord := range p.RunRecords {
		p.Memos[runRecord.FormulaHID] = runRecord.HID
	}
	// Rebuild `Whereabouts` by whitelisting prev values still ref'd by `Tags`.
	for _, release := range p.Tags {
		whereabout, ok := oldWhereabouts[release.Ware]
		if !ok {
			continue // fine; not everything is required to have this metadata.
		}
		p.Whereabouts[release.Ware] = whereabout
	}
}
