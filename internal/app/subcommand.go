package app

import (
	"flag"
)

var SupportedSubCommands = SubCommands{
	SubCommandReverse,
	SubCommandResize,
	SubCommandTile,
	SubCommandTrim,
	SubCommandGrayscale,
	SubCommandAddstring,
	SubCommandPng,
	SubCommandJpeg,
	SubCommandGif,
}

var SubCommandPng = &SubCommand{
	Name:            "png",
	Usage:           "file convert to png",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{},
}

var SubCommandJpeg = &SubCommand{
	Name:            "jpeg",
	Usage:           "file convert to jpeg",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{},
}

var SubCommandGif = &SubCommand{
	Name:            "gif",
	Usage:           "file convert to gif",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{},
}

var SubCommandReverse = &SubCommand{
	Name:            "reverse",
	Usage:           "reverse image",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{OptionVertical},
}

var SubCommandResize = &SubCommand{
	Name:            "resize",
	Usage:           "resize image",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{OptionWidth, OptionHeight, OptionRatio},
}

var SubCommandTile = &SubCommand{
	Name:            "tile",
	Usage:           "lay down images with x * y",
	RequiredOptions: []Option{OptionX, OptionY},
	OptionalOptions: []Option{},
}

var SubCommandTrim = &SubCommand{
	Name:            "trim",
	Usage:           "trim image",
	RequiredOptions: []Option{OptionLeft, OptionTop, OptionWidth, OptionHeight},
	OptionalOptions: []Option{},
}

var SubCommandGrayscale = &SubCommand{
	Name:            "grayscale",
	Usage:           "change image color to grayscale",
	RequiredOptions: []Option{},
	OptionalOptions: []Option{},
}

var SubCommandAddstring = &SubCommand{
	Name:            "addstring",
	Usage:           "add string on image",
	RequiredOptions: []Option{OptionText},
	OptionalOptions: []Option{OptionTtf, OptionSize, OptionTop, OptionLeft, OptionColor},
}

// SubCommand imgedit subcommand
type SubCommand struct {
	Name            string
	Usage           string
	RequiredOptions []Option
	OptionalOptions []Option
}

// ValidOption check the validity of options
func (s *SubCommand) ValidOption() bool {
	var requiredCount int
	var optional = true
	flag.Visit(func(f *flag.Flag) {
		for _, v := range s.RequiredOptions {
			if f.Name == v.Name() {
				requiredCount++
				return
			}
		}
		for _, v := range s.OptionalOptions {
			if f.Name == v.Name() {
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
