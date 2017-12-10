package main

import (
	"errors"
	"strings"

	"github.com/ElecProg/LamCalc"
)

func parseStatement(stmnt string) (cLCStatement, error) {
	// Some clean up
	stmnt = strings.TrimSpace(stmnt)

	if len(stmnt) == 0 || strings.HasPrefix(stmnt, "--") {
		// Empty line or comment
		return cLCStatement{command: "none"}, nil
	}

	switch strings.Fields(stmnt)[0] {
	case "exit":
		return cLCStatement{command: "exit"}, nil

	case "clear":
		return cLCStatement{command: "clear"}, nil

	case "info":
		return cLCStatement{command: "info"}, nil

	case "let":
		stmnt = strings.TrimPrefix(stmnt, "let")
		splitStmnt := strings.SplitAfter(stmnt, "=")

		if len(splitStmnt) < 2 {
			return cLCStatement{}, errors.New("no expression in let operation")
		}

		varname := strings.TrimSpace(strings.TrimSuffix(splitStmnt[0], "="))

		if strings.HasPrefix(varname, "\\") || strings.Contains(varname, " ") {
			return cLCStatement{}, errors.New("invalid variable name '" + varname + "' in let operation")
		}

		expression, err := LamCalc.ParseString(splitStmnt[1], globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "let",
			parameters: []interface{}{varname, expression},
		}, nil

	case "wlet":
		stmnt = strings.TrimPrefix(stmnt, "wlet")
		splitStmnt := strings.SplitAfter(stmnt, "=")

		if len(splitStmnt) < 2 {
			return cLCStatement{}, errors.New("no expression in wlet operation")
		}

		varname := strings.TrimSpace(strings.TrimSuffix(splitStmnt[0], "="))

		if strings.HasPrefix(varname, "\\") || strings.Contains(varname, " ") {
			return cLCStatement{}, errors.New("invalid variable name '" + varname + "' in wlet operation")
		}

		expression, err := LamCalc.ParseString(splitStmnt[1], globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "wlet",
			parameters: []interface{}{varname, expression},
		}, nil

	case "fold":
		stmnt = strings.TrimPrefix(stmnt, "fold")
		splitStmnt := strings.SplitAfter(stmnt, "into")

		if len(splitStmnt) < 2 {
			return cLCStatement{}, errors.New("no targets in fold operation")
		}

		expression, err := LamCalc.ParseString(strings.TrimSuffix(splitStmnt[0], "into"), globals)
		vars := strings.Fields(splitStmnt[1])

		if len(vars) == 0 {
			return cLCStatement{}, errors.New("no targets in fold operation")

		} else if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "fold",
			parameters: []interface{}{expression, vars},
		}, nil

	case "load":
		fields := strings.Fields(stmnt)

		if len(fields) > 1 {
			return cLCStatement{
				command:    "load",
				parameters: []interface{}{fields[1:]},
			}, nil
		}

		return cLCStatement{}, errors.New("no files listed to load")

	case "weak":
		stmnt = strings.TrimPrefix(stmnt, "weak")
		expression, err := LamCalc.ParseString(stmnt, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "weak",
			parameters: []interface{}{expression},
		}, nil

	default:
		expression, err := LamCalc.ParseString(stmnt, globals)

		if err != nil {
			return cLCStatement{}, err
		}

		return cLCStatement{
			command:    "show",
			parameters: []interface{}{expression},
		}, nil
	}
}
