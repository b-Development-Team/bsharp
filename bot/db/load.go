package db

import (
	"bufio"
	"bytes"
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

type dataOp struct {
	Key   string
	Value string
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		out.Write(line)
		if !isPrefix {
			return out.Bytes(), nil
		}
	}
}

func (d *Data) loadData() error {
	err := os.MkdirAll(filepath.Join(d.path, "data"), os.ModePerm)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(filepath.Join(d.path, "data"))
	if err != nil {
		return err
	}

	var op dataOp
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			// Load program
			f, err := os.OpenFile(filepath.Join(d.path, "data", file.Name()), os.O_RDWR, os.ModePerm)
			if err != nil {
				return err
			}

			// Read
			dat := make(map[string]string)
			reader := bufio.NewReader(f)
			size := 0
			for {
				line, err := readLine(reader)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}

				err = json.Unmarshal(line, &op)
				if err != nil {
					return err
				}

				dat[op.Key] = op.Value
				size += len(op.Key) + len(op.Value)
			}

			// Save
			name := strings.TrimSuffix(file.Name(), ".json")
			d.data[name] = dat
			d.dataFiles[name] = f
			d.datasize += size
		}
	}
	return nil
}
