package LamCalc

import (
	"strconv"
	"strings"
)

func intToLetter(num int) string {
	if num < 3 {
		// x, y, z
		return string(rune(120 + num))

	} else if num < 6 {
		// u, v, w
		return string(rune(117 - 3 + num))

	} else if num < 26 {
		// a, b, c
		return string(rune(num - 3 + 97))
	}

	return strconv.Itoa(num - 26)
}

// String returns the Lambda Expression as a string
func (lx LamExpr) String() string {
	return lx.deDebruijn([]string{}, 0)
}

func (lx LamExpr) deDebruijn(boundLetters []string, nextletter int) string {
	result := ""

	for _, part := range lx {
		switch part := part.(type) {
		case int:
			if part < int(len(boundLetters)) && boundLetters[part] != "" {
				result += boundLetters[part] + " "
			} else {
				newLetter := intToLetter(nextletter)
				nextletter++

				for i := int(len(boundLetters)); i < part; i++ {
					boundLetters = append(boundLetters, "")
				}

				boundLetters = append(boundLetters, newLetter)
				result += newLetter + " "
			}

		case LamFunc:
			if len(lx) == 1 {
				result += part.deDebruijn(boundLetters, nextletter)
			} else {
				result = strings.TrimSuffix(result, " ") + "(" + part.deDebruijn(boundLetters, nextletter) + ") "
			}

		case LamExpr:
			if len(lx) == 1 {
				result += part.deDebruijn(boundLetters, nextletter)
			} else {
				result = strings.TrimSuffix(result, " ") + "(" + part.deDebruijn(boundLetters, nextletter) + ") "
			}

		default:
			panic("invalid type in LamExpr")
		}
	}

	return strings.TrimSuffix(result, " ")
}

// String returns the Lambda Function as a string
func (lf LamFunc) String() string {
	return lf.deDebruijn([]string{}, 0)
}

func (lf LamFunc) deDebruijn(boundLetters []string, nextletter int) string {
	// First make the first character undefined (for now)
	newLetter := intToLetter(nextletter)
	nextletter++

	boundLetters = append([]string{newLetter}, boundLetters...)
	result := "L" + newLetter + "."

	lx := LamExpr(lf)
	result += lx.deDebruijn(boundLetters, nextletter)

	return result
}
