package mal

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

//CoreNS contains builtin functions for mal
var CoreNS = map[*Symbol]*Function{
	&Symbol{Value: "+"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Number{Value: a.Value + b.Value}, nil
	}},
	&Symbol{Value: "-"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Number{Value: a.Value - b.Value}, nil
	}},
	&Symbol{Value: "*"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Number{Value: a.Value * b.Value}, nil
	}},
	&Symbol{Value: "/"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Number{Value: a.Value / b.Value}, nil
	}},
	//take the parameters and return them as a list.
	&Symbol{Value: "list"}: &Function{Fn: func(args ...Type) (Type, error) {
		return &List{Value: args}, nil
	}},
	//return true if the first parameter is a list, false otherwise.
	&Symbol{Value: "list?"}: &Function{Fn: func(args ...Type) (Type, error) {
		_, ok := args[0].(*List)
		return &Boolean{Value: ok}, nil
	}},
	//treat the first parameter as a list and return true if the list is empty and false if it contains any elements.
	&Symbol{Value: "empty?"}: &Function{Fn: func(args ...Type) (Type, error) {
		if lst, ok := args[0].(*List); ok {
			if len(lst.Value) == 0 {
				return &Boolean{Value: true}, nil
			}
		}
		return &Boolean{Value: false}, nil
	}},
	// treat the first parameter as a list and return the number of elements that it contains.
	&Symbol{Value: "count"}: &Function{Fn: func(args ...Type) (Type, error) {
		if lst, ok := args[0].(*List); ok {
			return &Number{Value: float64(len(lst.Value))}, nil
		}
		return &Number{Value: 0}, nil
	}},

	&Symbol{Value: "<"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Boolean{Value: a.Value < b.Value}, nil
	}},
	&Symbol{Value: ">"}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Boolean{Value: a.Value > b.Value}, nil
	}},

	&Symbol{Value: "<="}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Boolean{Value: a.Value <= b.Value}, nil
	}},

	&Symbol{Value: ">="}: &Function{Fn: func(args ...Type) (Type, error) {
		a, _ := args[0].(*Number)
		b, _ := args[1].(*Number)
		return &Boolean{Value: a.Value >= b.Value}, nil
	}},

	&Symbol{Value: "pr-str"}: &Function{Fn: func(args ...Type) (Type, error) {
		var sb strings.Builder
		for i, v := range args {
			sb.WriteString(PrString(v, true))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}
		return &String{Value: sb.String()}, nil
	}},
	&Symbol{Value: "str"}: &Function{Fn: func(args ...Type) (Type, error) {
		var sb strings.Builder
		for _, v := range args {
			sb.WriteString(PrString(v, false))
		}
		return &String{Value: sb.String()}, nil
	}},
	&Symbol{Value: "prn"}: &Function{Fn: func(args ...Type) (Type, error) {
		var sb strings.Builder
		for i, v := range args {
			sb.WriteString(PrString(v, true))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}
		fmt.Println(sb.String())
		return &Nil{}, nil
	}},
	&Symbol{Value: "println"}: &Function{Fn: func(args ...Type) (Type, error) {
		var sb strings.Builder
		for i, v := range args {
			sb.WriteString(PrString(v, false))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}
		fmt.Println(sb.String())
		return &Nil{}, nil
	}},

	&Symbol{Value: "read-string"}: &Function{Fn: func(args ...Type) (Type, error) {
		v, _ := args[0].(*String)
		return ReadStr(v.Value)
	}},

	&Symbol{Value: "slurp"}: &Function{Fn: func(args ...Type) (Type, error) {
		filename, _ := args[0].(*String)
		dat, err := ioutil.ReadFile(filename.Value)
		if err != nil {
			return nil, err
		}
		return &String{Value: string(dat)}, nil
	}},

	// compare the first two parameters and return true if they are the same type and
	// contain the same value. In the case of equal length lists, each element of the
	// list should be compared for equality and if they are the same return true, otherwise false.
	// if we use an anonymous function here, we can't recurse, but we need to recurse to compare lists
	// so we define this function at the bottom of the file and refer to it by name here
	&Symbol{Value: "="}: &Function{Fn: compareFunc},
}

func compareFunc(args ...Type) (Type, error) {
	if reflect.TypeOf(args[0]) != reflect.TypeOf(args[1]) {
		return &Boolean{Value: false}, nil
	}
	switch v := args[0].(type) {
	case *Symbol:
		v2, _ := args[1].(*Symbol)
		return &Boolean{Value: v.Value == v2.Value}, nil
	case *Number:
		v2, _ := args[1].(*Number)
		return &Boolean{Value: v.Value == v2.Value}, nil
	case *List:
		v2, _ := args[1].(*List)
		if len(v.Value) != len(v2.Value) {
			return &Boolean{Value: false}, nil
		}
		for i := range v.Value {
			r, _ := compareFunc(v.Value[i], v2.Value[i])
			rbool, _ := r.(*Boolean)
			if rbool.Value == false {
				return &Boolean{Value: false}, nil
			}
		}
		return &Boolean{Value: true}, nil
	case *Boolean:
		v2, _ := args[1].(*Boolean)
		return &Boolean{Value: v.Value == v2.Value}, nil
	case *Nil:
		return &Boolean{Value: true}, nil
	case *Function:
		return &Boolean{Value: false}, nil // Go cant == functions, false seems to make the most sense
	case *String:
		v2, _ := args[1].(*String)
		return &Boolean{Value: v.Value == v2.Value}, nil

	default:
		return nil, fmt.Errorf("No equals operation implemented for type: %T", v)
	}
}
