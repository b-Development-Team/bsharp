package db

import (
	"fmt"
	"time"
)

func (d *Data) GetProgram(id string) (*Program, Resp) {
	d.RLock()
	defer d.RUnlock()

	p, exists := d.Programs[id]
	if !exists {
		return nil, Resp{Msg: fmt.Sprintf("Program `%s` doesn't exist!", id), Suc: false}
	}
	return p, Resp{Msg: "", Suc: true}
}

func (d *Data) NewProgram(id, name, creator string) *Program {
	return &Program{
		ID:       id,
		Name:     name,
		Creator:  creator,
		Uses:     0,
		Created:  time.Now(),
		LastUsed: time.Now(),
		Comment:  "",
		Image:    "",
	}
}

func (d *Data) GetProgramByName(name string) (string, Resp) {
	d.RLock()
	defer d.RUnlock()

	id, exists := d.names[name]
	if !exists {
		return "", Resp{Msg: fmt.Sprintf("Program **%s** doesn't exist!", name), Suc: false}
	}
	return id, Resp{Msg: "", Suc: true}
}

func (d *Data) GetSource(id string) (string, Resp) {
	d.RLock()
	defer d.RUnlock()

	src, exists := d.source[id]
	if !exists {
		return "", Resp{Msg: fmt.Sprintf("Source for `%s` doesn't exist!", id), Suc: false}
	}
	return src, Resp{Msg: "", Suc: true}
}
