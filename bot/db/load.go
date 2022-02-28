package db

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (d *Data) loadPrograms() error {
	err := os.MkdirAll(filepath.Join(d.path, "programs"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.path, "programs"))
	if err != nil {
		return err
	}
	var dat *Program
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			// Load program
			f, err := os.OpenFile(filepath.Join(d.path, "programs", file.Name()), os.O_RDWR, os.ModePerm)
			if err != nil {
				return err
			}
			dec := json.NewDecoder(f)
			err = dec.Decode(&dat)
			if err != nil {
				return err
			}
			d.Programs[dat.ID] = dat
			d.programFiles[dat.ID] = f
			d.names[dat.Name] = dat.ID
			dat = nil
		}
	}
	return nil
}

func (d *Data) loadSource() error {
	err := os.MkdirAll(filepath.Join(d.path, "source"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.path, "source"))
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".bsp") {
			// Load program
			f, err := os.OpenFile(filepath.Join(d.path, "source", file.Name()), os.O_RDWR, os.ModePerm)
			if err != nil {
				return err
			}
			dat, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(file.Name(), ".bsp")
			d.source[name] = string(dat)
			d.sourceFiles[name] = f
		}
	}
	return nil
}
