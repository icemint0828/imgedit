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
	OptionalOptions: []*Option{OptionWidth, OptionHeight},
}

var SubCommandTrim = &SubCommand{
	Name:            "trim",
	Usage:           "trim image",
	RequiredOptions: []*Option{},
	OptionalOptions: []*Option{OptionLeft, OptionRight, OptionBottom, OptionTop},
}

var SubCommandGrayscale = &SubCommand{
	Name:            "grayscale",
	Usage:           "change image color to grayscale",
	RequiredOptions: []*Option{},
	OptionalOptions: []*Option{},
}

type SubCommand struct {
	Name            string
	Usage           string
	RequiredOptions []*Option
	OptionalOptions []*Option
}

var OptionVertical = &Option{Name: "vertical", Usage: "direction for reverse. default horizon."}
var OptionWidth = &Option{Name: "width", Usage: "width for resize. default 100px."}
var OptionHeight = &Option{Name: "height", Usage: "height for resize. default 100px."}
var OptionLeft = &Option{Name: "left", Usage: "left for trim. default 0px."}
var OptionRight = &Option{Name: "right", Usage: "right for trim. default 100px."}
var OptionBottom = &Option{Name: "bottom", Usage: "bottom for trim. default 0px."}
var OptionTop = &Option{Name: "top", Usage: "top for trim. default 100px."}

type Option struct {
	Name  string
	Usage string
}

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
