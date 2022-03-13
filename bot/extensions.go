package bot

import (
	"errors"
	"fmt"

	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type extensionCtx struct {
	Dat    *db.Data
	Tag    string
	CurrDB string

	// Info
	Author string
}

func NewExtensionCtx(id string, dat *db.Data, ctx *Ctx) *extensionCtx {
	extCtx := &extensionCtx{
		Dat:    dat,
		Tag:    id,
		CurrDB: id,
		Author: ctx.Author(),
	}
	return extCtx
}

var exts = []*ir.Extension{
	{
		Name:    "DB",
		Params:  []types.Type{types.IDENT, types.STRING, types.VARIADIC},
		RetType: types.STRING,
	},
	{
		Name:    "USERID",
		Params:  []types.Type{},
		RetType: types.STRING,
	},
}

func getExtensions(c *extensionCtx) []*interpreter.Extension {
	return []*interpreter.Extension{
		interpreter.NewExtension("DB", func(pars []interface{}) (interface{}, error) {
			// Get args
			op := pars[0].(string)
			args := make([]string, len(pars)-1)
			for ind, arg := range pars[1:] {
				args[ind] = arg.(string)
			}

			// Get result
			switch op {
			case "SET":
				if c.CurrDB != c.Tag {
					return nil, fmt.Errorf("db: cannot write to read-only database")
				}
				if len(args) != 2 {
					return nil, fmt.Errorf("db: SET requires 2 arguments (key, value)")
				}
				err := c.Dat.DataSet(c.CurrDB, args[0], args[1])
				if err != nil {
					return nil, err
				}
				return "", nil

			case "USE":
				if len(args) != 1 {
					return nil, fmt.Errorf("db: USE requires 1 argument (database name)")
				}
				c.Dat.RLock()
				_, exists := c.Dat.Programs[args[0]]
				c.Dat.RUnlock()
				if !exists {
					return nil, fmt.Errorf("db: DB `%s` does not exist", args[0])
				}
				c.CurrDB = args[0]
				return "", nil

			case "GET":
				if len(args) != 1 {
					return nil, fmt.Errorf("db: GET requires 1 argument (key)")
				}
				val, rsp := c.Dat.DataGet(c.CurrDB, args[0])
				if !rsp.Suc {
					return nil, errors.New(rsp.Msg)
				}
				return val, nil
			}

			return nil, fmt.Errorf("db: unknown operation %s", op)
		}, []types.Type{types.IDENT, types.STRING, types.VARIADIC}, types.STRING),
		interpreter.NewExtension("USERID", func(pars []interface{}) (interface{}, error) {
			return c.Author, nil
		}, []types.Type{}, types.STRING),
	}
}
