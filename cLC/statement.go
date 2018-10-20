package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ElecProg/lamcalc"
)

type cLCStatement struct {
	command    string
	parameters []interface{}
}

func executeStatement(stmnt cLCStatement) {
	cmd := stmnt.command

	switch cmd {
	case "none":
		// Nothing to do

	case "exit":
		os.Exit(0)

	case "help":
		showHelp()

	case "clear":
		var clearCmd *exec.Cmd

		if runtime.GOOS == "windows" {
			clearCmd = exec.Command("cmd", "/c", "cls")
		} else {
			clearCmd = exec.Command("clear")
		}

		clearCmd.Stdout = os.Stdout
		clearCmd.Run()

	case "info":
		showInfo()

	case "let", "wlet":
		var rs lamcalc.Term
		switch cmd {
		case "wlet":
			rs = stmnt.parameters[1].(lamcalc.Term).WHNF()

		default:
			rs = concurrentReduce(stmnt.parameters[1].(lamcalc.Term))
		}

		if rs == nil {
			// Something went wrong
			return
		}

		globals[stmnt.parameters[0].(string)] = rs

	case "match":
		rs := concurrentReduce(stmnt.parameters[0].(lamcalc.Term))

		if rs == nil {
			// Something went wrong
			return
		}

		fmt.Print("\n" + stmnt.parameters[0].(lamcalc.Term).String() + " =\n")

		couldFold := false
		for _, gvar := range stmnt.parameters[1].([]string) {
			global, ok := globals[gvar]

			if !ok {
				global = lamcalc.Free(gvar)
			}

			if global.AlphaEquivalent(rs) {
				fmt.Print(rs.String() + " =\n\n")
				fmt.Print("    " + gvar + "\n\n")
				couldFold = true
				break
			}
		}

		if !couldFold {
			fmt.Print("\n    " + rs.String() + "\n\n")
		}

	case "load":
		loadFiles(stmnt.parameters[0].([]string))

	case "show", "weak":
		var rs lamcalc.Term

		switch cmd {
		case "weak":
			rs = stmnt.parameters[0].(lamcalc.Term).WHNF()

		default:
			rs = concurrentReduce(stmnt.parameters[0].(lamcalc.Term))
		}

		if rs == nil {
			// Something went wrong
			return
		}

		fmt.Print("\n" + stmnt.parameters[0].(lamcalc.Term).String() + " =\n\n")
		fmt.Print("    " + rs.String() + "\n\n")
	}
}
