package model

import (
	"fmt"

	rdef "go.polydawn.net/repeatr/api/def"
)

type Project struct {
	Names      map[string]ReleaseRecord   // map name->{ware,backstory}
	RunRecords map[string]*rdef.RunRecord // map rrhid->rr
	Memos      map[string]string          // index frmhid->rrhid
}

type ReleaseRecord struct {
	Ware         rdef.Ware
	RunRecordHID string // blank if a name was manual
}

func (p *Project) Init() {
	p.Names = make(map[string]ReleaseRecord)
	p.RunRecords = make(map[string]*rdef.RunRecord)
	p.Memos = make(map[string]string)
}

func (p *Project) PutManualName(name string, ware rdef.Ware) {
	_, hadPrev := p.Names[name]
	p.Names[name] = ReleaseRecord{ware, ""}
	if hadPrev {
		p.retainFilter()
	}
}

func (p *Project) DeleteName(name string) {
	_, hadPrev := p.Names[name]
	if hadPrev {
		delete(p.Names, name)
		p.retainFilter()
	}
}

func (p *Project) PutEval(rr *rdef.RunRecord) {
	var savedAny bool
	for name, value := range rr.Results {
		if (name[0] >= 'a' && name[0] <= 'z') ||
			(name[0] >= 'A' && name[0] <= 'Z') {
			savedAny = true
		}
		p.Names[name] = ReleaseRecord{value.Ware, rr.HID}

	}
	// If no interesting results were saved, bail;
	//  we don't need to save this record, nor potentially evict old ones.
	// REVIEW: this means we won't memoize evals that didn't have named results.
	if savedAny == false {
		return
	}
	p.RunRecords[rr.HID] = rr
	p.Memos[rr.FormulaHID] = rr.HID
	p.retainFilter()
}

func (p *Project) retainFilter() {
	// "Sweep".  (The `Names` map is the marks.)
	oldRunRecords := p.RunRecords
	p.RunRecords = make(map[string]*rdef.RunRecord)
	p.Memos = make(map[string]string)
	// Rebuild `RunRecords` by whitelisting prev values still ref'd by `Names`.
	for name, release := range p.Names {
		if release.RunRecordHID == "" {
			continue // skip.  it's just a fiat release; doesn't ref anything.
		}
		runRecord, ok := oldRunRecords[release.RunRecordHID]
		if !ok {
			panic(fmt.Errorf("db integrity violation: dangling runrecord -- release %q points to %q", name, release.RunRecordHID))
		}
		p.RunRecords[release.RunRecordHID] = runRecord
	}
	// Rebuild `Memos` index from `RunRecords`.
	for _, runRecord := range p.RunRecords {
		p.Memos[runRecord.FormulaHID] = runRecord.HID
	}
}
