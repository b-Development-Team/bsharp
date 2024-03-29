package bot

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
	"github.com/Nv7-Github/sevcord"
)

type extensionCtx struct {
	Dat    *db.Data
	Tag    string
	CurrDB string
	Stdout *ctxWriter

	// Info
	Author      string
	Multiplayer bool
}

func NewExtensionCtx(id string, dat *db.Data, ctx sevcord.Ctx) *extensionCtx {
	extCtx := &extensionCtx{
		Dat:    dat,
		Tag:    id,
		CurrDB: id,
		Author: ctx.User().ID,
	}
	return extCtx
}

var Exts = []*ir.Extension{
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
	{
		Name:    "INPUT",
		Params:  []types.Type{types.STRING},
		RetType: types.STRING,
	},
	{
		Name:    "MULTIPLAYER", // allows anyone to interact with the buttons and input, or disables
		Params:  []types.Type{types.BOOL},
		RetType: types.NULL,
	},
	{
		Name:    "BUTTON",
		Params:  []types.Type{types.STRING, types.STRING}, // text, id
		RetType: types.NewMapType(types.STRING, types.STRING),
	},
	{
		Name:    "BUTTONS",
		Params:  []types.Type{types.NewArrayType(types.NewArrayType(types.NewMapType(types.STRING, types.STRING)))}, // [][]btn
		RetType: types.STRING,                                                                                       // id of button pressed
	},
	{
		Name:    "DISABLE",
		Params:  []types.Type{types.NewMapType(types.STRING, types.STRING)},
		RetType: types.NewMapType(types.STRING, types.STRING),
	},
	{
		Name:    "ENABLE",
		Params:  []types.Type{types.NewMapType(types.STRING, types.STRING)},
		RetType: types.NewMapType(types.STRING, types.STRING),
	},
	{
		Name:    "COLOR",
		Params:  []types.Type{types.NewMapType(types.STRING, types.STRING), types.INT},
		RetType: types.NewMapType(types.STRING, types.STRING),
	},
}

func getExtensions(c *extensionCtx) []*interpreter.Extension {
	return []*interpreter.Extension{
		interpreter.NewExtension("DB", func(pars []any) (any, error) {
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

			case "EXISTS":
				if len(args) != 1 {
					return nil, fmt.Errorf("db: EXISTS requires 1 argument (key)")
				}

				val := c.Dat.DataExists(c.CurrDB, args[0])
				if val {
					return "true", nil
				}
				return "false", nil
			}

			return nil, fmt.Errorf("db: unknown operation %s", op)
		}, []types.Type{types.IDENT, types.STRING, types.VARIADIC}, types.STRING),
		interpreter.NewExtension("USERID", func(pars []any) (any, error) {
			return c.Author, nil
		}, []types.Type{}, types.STRING),
		interpreter.NewExtension("MULTIPLAYER", func(pars []any) (any, error) {
			c.Multiplayer = pars[0].(bool)
			return nil, nil
		}, []types.Type{}, types.STRING),
		interpreter.NewExtension("INPUT", func(pars []any) (any, error) {
			out := make(chan string)
			prompt := pars[0].(string)
			if len(prompt) > 45 {
				return nil, errors.New("input: prompt must be 45 or fewer characters")
			}
			emb := sevcord.NewEmbedBuilder(prompt).Description("Press the button bellow to respond.").Color(16776960)
			c.Stdout.Edit(sevcord.EmbedResponse(emb).ComponentRow(&sevcord.Button{
				Label: "Respond",
				Style: sevcord.ButtonStyleSuccess,
				Handler: func(ctx sevcord.Ctx) {
					if !c.Multiplayer && ctx.User().ID != c.Stdout.User().ID {
						return
					}

					ctx.Modal(&sevcord.Modal{
						Title: "Respond",
						Inputs: []sevcord.ModalInput{
							{
								Label:       prompt,
								Style:       sevcord.ModalInputStyleParagraph,
								Placeholder: `Your response...`,
								Required:    true,
								MaxLength:   4000,
								MinLength:   1,
							},
						},
						Handler: func(ctx sevcord.Ctx, args []string) {
							v := args[0]
							c.Stdout.Ctx = ctx
							out <- v
						},
					})
				},
			}))

			select {
			case rsp := <-out:
				return rsp, nil

			case <-time.After(time.Second * 30):
				return "", errors.New("user took more than 30 seconds to respond")
			}
		}, []types.Type{types.STRING}, types.STRING),
		interpreter.NewExtension("BUTTON", func(pars []any) (any, error) {
			return map[string]any{
				"id":  pars[1].(string),
				"txt": pars[0].(string),
				"col": "Primary",
				"dis": "false",
			}, nil
		}, []types.Type{types.STRING, types.STRING}, types.NewMapType(types.STRING, types.STRING)),
		interpreter.NewExtension("DISABLE", func(pars []any) (any, error) {
			pars[0].(map[string]any)["dis"] = "true"
			return pars[0], nil
		}, []types.Type{types.NewMapType(types.STRING, types.STRING)}, types.NewMapType(types.STRING, types.STRING)),
		interpreter.NewExtension("ENABLE", func(pars []any) (any, error) {
			pars[0].(map[string]any)["dis"] = "false"
			return pars[0], nil
		}, []types.Type{types.NewMapType(types.STRING, types.STRING)}, types.NewMapType(types.STRING, types.STRING)),
		interpreter.NewExtension("COLOR", func(pars []any) (any, error) {
			cols := map[int]string{
				1: "Primary",
				2: "Secondary",
				3: "Danger",
				4: "Success",
			}
			col, exists := cols[pars[1].(int)]
			if !exists {
				return nil, errors.New("invalid color")
			}
			pars[0].(map[string]any)["col"] = col
			return pars[0], nil
		}, []types.Type{types.NewMapType(types.STRING, types.STRING), types.INT}, types.NewMapType(types.STRING, types.STRING)),
		interpreter.NewExtension("BUTTONS", func(pars []any) (any, error) {
			// Build actions row
			r := *(pars[0].(*[]any))
			rows := make([][]sevcord.Component, len(r))

			var cols = map[string]sevcord.ButtonStyle{
				"Primary":   sevcord.ButtonStylePrimary,
				"Secondary": sevcord.ButtonStyleSecondary,
				"Danger":    sevcord.ButtonStyleDanger,
				"Success":   sevcord.ButtonStyleSuccess,
			}

			out := make(chan string)
			for i, val := range r {
				v := *val.(*[]any)
				r := make([]sevcord.Component, len(v))
				for j, btn := range v {
					b := btn.(map[string]any)
					col, ok := b["col"].(string)
					if !ok {
						return nil, errors.New("invalid color")
					}
					s, exists := cols[col]
					if !exists {
						return nil, errors.New("invalid color")
					}

					val := b["id"].(string)
					handler := func(ctx sevcord.Ctx) {
						if !c.Multiplayer && ctx.User().ID != c.Stdout.User().ID {
							return
						}

						c.Stdout.Ctx = ctx
						c.Stdout.Flush()
						out <- val
					}

					txt := b["txt"].(string)
					if len([]rune(txt)) == 1 && []rune(txt)[0] > unicode.MaxASCII {
						r[j] = &sevcord.Button{
							Emoji:    sevcord.ComponentEmojiDefault([]rune(txt)[0]),
							Style:    s,
							Disabled: b["dis"] == "true",
							Handler:  handler,
						}
					} else {
						r[j] = &sevcord.Button{
							Label:    txt,
							Style:    s,
							Disabled: b["dis"] == "true",
							Handler:  handler,
						}
					}
				}

				rows[i] = r
			}

			// Send
			c.Stdout.cmps = rows
			err := c.Stdout.Flush()
			if err != nil {
				return nil, err
			}

			select {
			case rsp := <-out:
				return rsp, nil

			case <-time.After(time.Second * 30):
				return "", errors.New("user took more than 30 seconds to respond")
			}
		}, []types.Type{types.NewArrayType(types.NewArrayType(types.NewMapType(types.STRING, types.STRING)))}, types.STRING),
	}
}
