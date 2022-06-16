package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/icemint0828/imgedit/internal/app"
)

const (
	Version = "0.0.1"
)

func main() {
	flag.Bool(app.OptionVertical.Name, false, app.OptionVertical.Usage)
	flag.Uint(app.OptionWidth.Name, 100, app.OptionWidth.Usage)
	flag.Uint(app.OptionHeight.Name, 100, app.OptionHeight.Usage)
	flag.Float64(app.OptionRatio.Name, 0, app.OptionRatio.Usage)
	flag.Uint(app.OptionLeft.Name, 0, app.OptionLeft.Usage)
	flag.Uint(app.OptionRight.Name, 100, app.OptionRight.Usage)
	flag.Uint(app.OptionBottom.Name, 0, app.OptionBottom.Usage)
	flag.Uint(app.OptionTop.Name, 100, app.OptionTop.Usage)
	flag.CommandLine.Usage = usage
	permuteArgs(os.Args[1:])
	flag.Parse()

	// validation for flag and args
	args := flag.Args()

	if len(args) != 2 {
		fmt.Println("Argument is missing.")
		usage()
		os.Exit(1)
	}
	subCommandName, imagePath := args[0], args[1]
	subCommand := app.SupportedSubCommands.FindSubCommand(subCommandName)
	if subCommand == nil {
		fmt.Printf("%s is not supported for subcommand.\n", subCommandName)
		usage()
		os.Exit(1)
	}

	if !subCommand.ValidOption() {
		fmt.Printf("%s is not valid for option.\n", subCommandName)
		usage()
		os.Exit(1)
	}

	if !exists(imagePath) {
		fmt.Printf("File does not exist : %s\n", imagePath)
		os.Exit(1)
	}

	if !supportedExtension(imagePath) {
		fmt.Printf("Extension is not supported : %s\n", imagePath)
		os.Exit(1)
	}

	// run application
	err := app.NewApp(subCommand, imagePath).Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func supportedExtension(filename string) bool {
	e := filepath.Ext(filename)
	for _, extension := range app.SupportedExtensions {
		if strings.ToLower(e) == extension {
			return true
		}
	}
	return false
}

func usage() {
	commandName := path.Base(os.Args[0])
	fmt.Printf("%s version %s\n\n", commandName, Version)
	fmt.Printf("Try running %s like:\n", commandName)
	fmt.Printf("%s <sub command> <image path> -<option> | for example:\n\n", commandName)
	fmt.Printf("%s reverse test.png -vertical\n", commandName)
	fmt.Printf("%s resize test.png -width 500 -height 500\n\n", commandName)
	fmt.Printf("[sub command]\n")
	for _, subCommand := range app.SupportedSubCommands {
		fmt.Printf("\n  %s : %s\n", subCommand.Name, subCommand.Usage)
		//fmt.Printf("    [required options]\n")
		//for _, option := range subCommand.RequiredOptions {
		//	fmt.Printf("      -%s : %s\n", option.Name, option.Usage)
		//}
		if len(subCommand.OptionalOptions) == 0 {
			continue
		}
		fmt.Printf("    [optional options]\n")
		for _, option := range subCommand.OptionalOptions {
			fmt.Printf("      -%s : %s\n", option.Name, option.Usage)
		}
	}
	fmt.Printf("\n[supported extensions]\n")
	fmt.Printf("%s\n", strings.Join(app.SupportedExtensions, "/"))
}

func permuteArgs(args []string) {
	var flagArgs []string
	var nonFlagArgs []string

	for i := 0; i < len(args); i++ {
		v := args[i]
		if v[0] == '-' {
			optionName := v[1:]
			switch optionName {
			case app.OptionHeight.Name, app.OptionWidth.Name, app.OptionRatio.Name, app.OptionLeft.Name, app.OptionRight.Name, app.OptionBottom.Name, app.OptionTop.Name:
				/* out of index */
				if len(args) <= i+1 {
					usage()
					os.Exit(1)
				}
				/* the next flag has come */
				optionVal := args[i+1]
				if optionVal[0] == '-' {
					usage()
					os.Exit(1)
				}
				flagArgs = append(flagArgs, args[i:i+2]...)
				i++
			default:
				flagArgs = append(flagArgs, args[i])
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, args[i])
		}
	}
	permutedArgs := append(flagArgs, nonFlagArgs...)

	/* replace args */
	for i := 0; i < len(args); i++ {
		args[i] = permutedArgs[i]
	}
}
