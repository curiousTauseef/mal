package mal

import (
	"strconv"
	"strings"
)

//PrString takes a MalType and returns a string representation
func PrString(ast Type) string {
	var sb strings.Builder

	switch v := ast.(type) {
	case *List:
		sb.WriteString("(")
		for i, vel := range v.value {
			sb.WriteString(printAtom(vel))
			if i < len(v.value)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(")")
	default:
		sb.WriteString(printAtom(v))

	}
	return sb.String()
}

func printAtom(atom Type) string {
	switch v := atom.(type) {
	case *Symbol:
		return v.value
	case *Number:
		// format the float as a string, g as parameter means:
		// see https://golang.org/pkg/strconv/#FormatFloat
		// g' ('e' for large exponents, 'f' otherwise)
		//  'e' (-d.dddde±dd, a decimal exponent)
		//  'f' (-ddd.dddd, no exponent)
		return strconv.FormatFloat(v.value, 'g', -1, 64)
	case *List:
		return PrString(v)

	default:
		return "<TO STRING NOT IMPLEMENTED>"
	}
}
