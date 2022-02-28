package main

import (
	"errors"
	"os"

	"github.com/Nv7-Github/bsharp/bot/db"
)

func main() {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		panic(err)
	}
	d, err := db.NewDB("data")
	if err != nil {
		panic(err)
	}

	dat, err := d.Get("epicgld")
	if err != nil {
		panic(err)
	}

	id, rsp := dat.GetProgramByName("Epic Program")
	if !rsp.Suc {
		prog := dat.NewProgram("epicprogram", "Epic Program", "userid")
		err = dat.SaveProgram(prog)
		if err != nil {
			panic(err)
		}
		id = "epicprogram"
	}
	prog, rsp := dat.GetProgram(id)
	if !rsp.Suc {
		panic(errors.New(rsp.Msg))
	}
	prog.Comment = "Very Epic Program"
	err = dat.SaveProgram(prog)
	if err != nil {
		panic(err)
	}
	src := `[PRINT "This is an epic program!"]`
	err = dat.SaveSource(id, src)
	if err != nil {
		panic(err)
	}
}
