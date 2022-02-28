package db

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Resp struct {
	Msg string
	Suc bool
}

type Program struct {
	ID          string
	Name        string
	Creator     string // User ID
	Uses        int
	Created     time.Time
	LastUsed    time.Time
	Description string
	Image       string
}

type Data struct {
	*sync.RWMutex
	names    map[string]string // map[name]id
	Programs map[string]*Program
	source   map[string]string
	Guild    string
	path     string

	programFiles map[string]*os.File
	sourceFiles  map[string]*os.File
}

func newData(path, guild string) *Data {
	return &Data{
		RWMutex: &sync.RWMutex{},
		path:    path,

		Programs:     make(map[string]*Program),
		Guild:        guild,
		programFiles: make(map[string]*os.File),
		sourceFiles:  make(map[string]*os.File),
		names:        make(map[string]string),
		source:       make(map[string]string),
	}
}

func (d *Data) Close() {
	for _, f := range d.programFiles {
		f.Close()
	}
	for _, f := range d.sourceFiles {
		f.Close()
	}
}

type DB struct {
	*sync.RWMutex
	datas map[string]*Data
	path  string
}

func (d *DB) Close() {
	for _, data := range d.datas {
		data.Close()
	}
}

func (d *DB) Get(gld string) (*Data, error) {
	d.RLock()
	data, exists := d.datas[gld]
	d.RUnlock()
	if !exists {
		// Create if not exists
		data = newData(filepath.Join(d.path, gld), gld)
		err := data.loadPrograms()
		if err != nil {
			return nil, err
		}
		err = data.loadSource()
		if err != nil {
			return nil, err
		}
		d.Lock()
		d.datas[gld] = data
		d.Unlock()
	}
	return data, nil
}

func newDB(path string) *DB {
	return &DB{
		RWMutex: &sync.RWMutex{},
		datas:   make(map[string]*Data),
		path:    path,
	}
}

func NewDB(path string) (*DB, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	db := newDB(path)
	folders, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, folder := range folders {
		if folder.IsDir() {
			gld := folder.Name()
			data := newData(filepath.Join(path, gld), gld)
			err := data.loadPrograms()
			if err != nil {
				return nil, err
			}
			err = data.loadSource()
			if err != nil {
				return nil, err
			}
			db.datas[gld] = data
		}
	}
	return db, nil
}
