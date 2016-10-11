package model

import (
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

}
