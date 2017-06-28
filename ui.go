package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

var ui *cli.ColoredUi

func initUi() error {
	ui = new(cli.ColoredUi)
	if ui == nil {
		fmt.Printf("error of ui\n")
		return errors.New("failed to new cli")
	}

	bui := new(cli.BasicUi)
	bui.Reader = os.Stdin
	bui.Writer = os.Stdout
	bui.ErrorWriter = os.Stderr

	ui.Ui = bui
	ui.OutputColor = cli.UiColorNone
	ui.InfoColor = cli.UiColorGreen
	ui.ErrorColor = cli.UiColorRed
	ui.WarnColor = cli.UiColorYellow

	return nil
}

func printUsage() {
	ui.Output("Usage:")
	ui.Output("  [h]elp          : print this message.")
	ui.Output("  [p]rint <idx>   : show regualtors.")
	ui.Output("  exit            : exit this program.")
}

func handleInput(input string) (exit bool) {
	exit = false

	cmdline := strings.Fields(input)
	if len(cmdline) == 0 {
		printUsage()
		return
	}

	switch cmdline[0] {
	case "exit":
		exit = true
	case "help", "h":
		printUsage()
	case "print", "p":
		index := -1
		if len(cmdline) >= 2 {
			index, _ = strconv.Atoi(cmdline[1])
		}
		cmdPrint(os.Stdout, regulators, index)
	default:
		printUsage()
	}

	return
}

func cmdPrint(w io.Writer, rs []*regulator, idx int) {
	if idx >= 0 {
		printRegulator(w, rs[idx], true)
		return
	}
	for _, r := range rs {
		printRegulator(w, r, false)
	}
}
