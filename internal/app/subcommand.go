package app

import (
	"flag"
)

var SupportedSubCommands = SubCommands{
	SubCommandReverse,
	SubCommandResize,
	SubCommandTrim,
	SubCommandGrayscale,
}

var SubCommandReverse = &SubCommand{
	Name:            "reverse",
	Usage:           "reverse image",
	RequiredOptions: []*Option{},
	OptionalOptions: []*Option{OptionVertical},
}

var SubCommandResize = &SubCommand{
	Name:            "resize",
	Usage:           "resize image",
	RequiredOptions: []*Option{},
	OptionalOptions: []*Option{OptionWidth, OptionHeight, OptionRatio},
}

var SubCommandTrim = &SubCommand{
	Name:            "trim",
	Usage:           "trim image",
	RequiredOptions: []*Option{OptionLeft, OptionTop, OptionWidth, OptionHeight},
	OptionalOptions: []*Option{},
}

var SubCommandGrayscale = &SubCommand{
	Name:            "grayscale",
	Usage:           "change image color to grayscale",
	RequiredOptions: []*Option{},
	OptionalOptions: []*Option{},
}

// SubCommand imgedit subcommand
type SubCommand struct {
	Name            string
	Usage           string
	RequiredOptions []*Option
	OptionalOptions []*Option
}

var OptionVertical = &Option{Name: "vertical", Usage: "direction for reverse. default horizon."}
var OptionWidth = &Option{Name: "width", Usage: "width px."}
var OptionHeight = &Option{Name: "height", Usage: "height px."}
var OptionRatio = &Option{Name: "ratio", Usage: "ratio for resize. if ratio is set, width and height are ignored."}
var OptionTop = &Option{Name: "top", Usage: "start top point px of trim."}
var OptionLeft = &Option{Name: "left", Usage: "start left point px of trim."}

// Option for subcommands
type Option struct {
	Name  string
	Usage string
}

// ValidOption check the validity of options
func (s *SubCommand) ValidOption() bool {
	var requiredCount int
	var optional = true
	flag.Visit(func(f *flag.Flag) {
		for _, v := range s.RequiredOptions {
			if f.Name == v.Name {
				requiredCount++
				return
			}
		}

		for _, v := range s.OptionalOptions {
			if f.Name == v.Name {
				return
			}
		}
		optional = false
	})
	return requiredCount == len(s.RequiredOptions) && optional
}

type SubCommands []*SubCommand

func (s SubCommands) FindSubCommand(subCommandName string) *SubCommand {
	for _, v := range s {
		if v.Name == subCommandName {
			return v
		}
	}
	return nil
}
