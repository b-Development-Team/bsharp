package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (d *Data) SaveProgram(p *Program) error {
	d.Lock()
	defer d.Unlock()

	var f *os.File
	prog, exists := d.Programs[p.ID]
	if !exists {
		var err error
		f, err = os.Create(filepath.Join(d.path, "programs", p.ID+".json"))
		if err != nil {
			return err
		}
		d.programFiles[p.ID] = f
		d.names[p.Name] = p.ID
	} else {
		f = d.programFiles[p.ID]
		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}
		err = f.Truncate(0)
		if err != nil {
			return err
		}
		if prog.Name != p.Name { // Changed name?
			delete(d.names, prog.Name)
			d.names[p.Name] = p.ID
		}
	}

	// Save
	d.Programs[p.ID] = p
	enc := json.NewEncoder(f)
	return enc.Encode(p)
}

func (d *Data) SaveSource(id, source string) error {
	d.Lock()
	defer d.Unlock()

	var f *os.File
	_, exists := d.source[id]
	if !exists {
		var err error
		f, err = os.Create(filepath.Join(d.path, "source", id+".bsp"))
		if err != nil {
			return err
		}
		d.sourceFiles[id] = f
	} else {
		f = d.sourceFiles[id]
		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}
		err = f.Truncate(0)
		if err != nil {
			return err
		}
	}

	// Save
	d.source[id] = source
	_, err := f.WriteString(source)
	return err
}

func (d *Data) DataSet(id string, key, val string) error {
	d.Lock()
	defer d.Unlock()

	if len(key) > MaxKeySize {
		return errors.New("db: key exceeds max size (256)")
	}

	// Get data
	v, exists := d.data[id]
	if !exists {
		v = make(map[string]string)
		d.data[id] = v
	}

	// Get file
	f, exists := d.dataFiles[id]
	if !exists {
		var err error
		f, err = os.Create(filepath.Join(d.path, "data", id+".json"))
		if err != nil {
			return err
		}
		d.dataFiles[id] = f
	}

	// Calc size
	newsize := d.datasize
	e, exists := v[key]
	if !exists {
		newsize += len(key) + len(val)
	} else {
		newsize -= len(e)
		newsize += len(val)
	}

	if newsize > MaxSize {
		return fmt.Errorf("db: max size (1MB) exceeded")
	}

	// Write
	op, err := json.Marshal(dataOp{
		Key:   key,
		Value: val,
	})
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(op) + "\n")
	if err != nil {
		return err
	}

	// Save
	v[key] = val

	return nil
}
