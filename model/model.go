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

func (p *Project) retainFilter() {

}
