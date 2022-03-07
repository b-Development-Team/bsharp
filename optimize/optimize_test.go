package optimize

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
)

const code = `# Functions
[FUNC HELLO 
  [PRINT "Hello, World!"]
]

# Some basic Functions
[DEFINE a 1]
[PRINT [STRING [VAR a]]]
[PRINT [CONCAT "hello w" "orld"]]
[IF [COMPARE 6 != 4]
  [PRINT "6 is not 4"]
ELSE
  [PRINT "6 is 4, uh oh!"]
]
[INDEX "Hi!" 0]
[PRINT [STRING [LENGTH "Hello, World!"]]]

[DEFINE i 0]
[WHILE [COMPARE [VAR i] < 100]
  [PRINT [STRING [VAR i]]]
  [DEFINE i [MATH [VAR i] + 1]]
]

[PRINT [STRING [FLOAT 100]]]

[PRINT [STRING [RANDINT 1 100]]]
[PRINT [STRING [RANDOM 1.0 100.0]]]
[PRINT [STRING [CEIL 1.5]]]
[PRINT [STRING [FLOOR 1.5]]]
[PRINT [STRING [ROUND 1.5]]]

# Maps
[DEFINE b [MAKE MAP{STRING}STRING]]
[SET [VAR b] "Hello" "World"]
[PRINT [GET [VAR b] "Hello"]]

# Arrays
[DEFINE c [ARRAY 1 2 3]]
[PRINT [STRING [INDEX [VAR c] 0]]]
[APPEND [VAR c] 4]

# Switch-case
[SWITCH "Hello"
  [CASE "World"
    [PRINT "The value was world!"]
  ]

  [CASE "Hello"
    [PRINT "The value was hello!"]
  ]

  [DEFAULT
    [PRINT "Unknown value!"]
  ]
]

# First-class Functions
[HELLO] # Call function with alias
[CALL [FN HELLO] ""] # Call function without alias
[DEFINE fns [MAKE MAP{STRING}FUNC{}NULL]] # Make a dictionary of functions
[SET [VAR fns] "hello" [FN HELLO]]
[CALL [GET [VAR fns] "hello"] ""] # Call function from dictionary
`

type fs struct{}

func (*fs) Parse(f string) (*parser.Parser, error) {
	return nil, errors.New("not implemented")
}

func TestOptimizer(t *testing.T) {
	stream := tokens.NewStream("main.bsp", code)
	tok := tokens.NewTokenizer(stream)
	err := tok.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	parse := parser.NewParser(tok)
	err = parse.Parse()
	if err != nil {
		t.Fatal(err)
	}
	cnf := ir.CodeConfig{
		Indent: 4,
	}
	ir := ir.NewBuilder()
	err = ir.Build(parse, &fs{})
	if err != nil {
		t.Fatal(err)
	}

	i := ir.IR()
	fmt.Println("BEFORE: ")
	fmt.Println(i.Code(cnf))

	opt := NewOptimizer(i)
	out := opt.Optimize()

	fmt.Println("AFTER: ")
	fmt.Println(out.Code(cnf))
}
