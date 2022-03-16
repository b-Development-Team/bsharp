package std

import (
	"embed"
	"io"
)

//go:embed *.bsp
var fs embed.FS

var Std = make(map[string]string)

func init() {
	dirs, err := fs.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, val := range dirs {
		f, err := fs.Open(val.Name())
		if err != nil {
			panic(err)
		}
		v, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		Std[val.Name()] = string(v)
		f.Close()
	}
}
