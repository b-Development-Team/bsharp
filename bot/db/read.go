package db

import (
	"fmt"
	"strings"
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

var validIDChars = "abcdefghijklmnopqrstuvwxyz_-0123456789"

const maxNameLength = 128

func (d *Data) NewProgram(id, name, creator string) (*Program, Resp) {
	// Validate id
	if id[0] == '_' {
		return nil, Resp{Msg: "IDs can't start with underscores!", Suc: false}
	}
	for _, char := range id {
		if !strings.Contains(validIDChars, string(char)) {
			return nil, Resp{Msg: fmt.Sprintf("Character \"%s\" isn't allowed in program IDs!", string(char)), Suc: false}
		}
	}
	if len([]rune(name)) > maxNameLength {
		return nil, Resp{Msg: "Program name must be under 128 characters!", Suc: false}
	}
	// Check if already exists
	_, rsp := d.GetProgram(id)
	if rsp.Suc {
		return nil, Resp{Msg: fmt.Sprintf("Tag with ID `%s` already exists!", id), Suc: false}
	}
	_, rsp = d.GetProgramByName(id)
	if rsp.Suc {
		return nil, Resp{Msg: fmt.Sprintf("Tag with name **%s** already exists!", name), Suc: false}
	}
	return &Program{
		ID:          id,
		Name:        name,
		Creator:     creator,
		Uses:        0,
		Created:     time.Now(),
		LastUsed:    time.Now(),
		Description: "None",
		Image:       "",
	}, Resp{Msg: "", Suc: true}
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

func (d *Data) DataGet(id string, key string) (string, Resp) {
	d.RLock()
	defer d.RUnlock()

	v, exists := d.data[id]
	if !exists {
		return "", Resp{Msg: fmt.Sprintf("Key **%s** doesn't exist!", key), Suc: false}
	}
	val, exists := v[key]
	if !exists {
		return "", Resp{Msg: fmt.Sprintf("Key **%s** doesn't exist!", key), Suc: false}
	}
	return val, Resp{Msg: "", Suc: true}
}
