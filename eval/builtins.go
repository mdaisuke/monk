package eval

import (
	"github.com/mdaisuke/monk/obj"
)

var builtins = map[string]*obj.Builtin{
	"len": &obj.Builtin{
		Fn: func(args ...obj.Obj) obj.Obj {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *obj.String:
				return &obj.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` is not supported, got %s", args[0].Type())
			}
		},
	},
}
