package app

import "flag"

var OptionVertical = &BoolOption{
	option: option{
		name:  "vertical",
		usage: "direction for reverse. default horizon.",
	},
	defaultVal: false,
}
var OptionWidth = &UintOption{
	option: option{
		name:  "width",
		usage: "width px.",
	},
	defaultVal: 0,
}
var OptionHeight = &UintOption{
	option: option{
		name:  "height",
		usage: "height px.",
	},
	defaultVal: 0,
}
var OptionRatio = &Float64Option{
	option: option{
		name:  "ratio",
		usage: "ratio for resize. if ratio is set, width and height are ignored.",
	},
	defaultVal: 0,
}
var OptionTop = &UintOption{
	option: option{
		name:  "top",
		usage: "start top point px.",
	},
	defaultVal: 0,
}
var OptionLeft = &UintOption{
	option: option{
		name:  "left",
		usage: "start left point px.",
	},
	defaultVal: 0,
}
var OptionX = &UintOption{
	option: option{
		name:  "x",
		usage: "x length for tile.",
	},
	defaultVal: 0,
}
var OptionY = &UintOption{
	option: option{
		name:  "y",
		usage: "y length for tile.",
	},
	defaultVal: 0,
}
var OptionText = &StringOption{
	option: option{
		name:  "text",
		usage: "text for addstring.",
	},
	defaultVal: "",
}
var OptionTtf = &StringOption{
	option: option{
		name:  "ttf",
		usage: "ttf file path for addstring.",
	},
	defaultVal: "",
}
var OptionSize = &UintOption{
	option: option{
		name:  "size",
		usage: "font size for addstring.",
	},
	defaultVal: 0,
}
var OptionColor = &StringOption{
	option: option{
		name:  "color",
		usage: "font color for addstring(back, white, red, blue, green).",
	},
	defaultVal: "",
}

// Option for subcommands
type Option interface {
	RegisterFlag()
	Name() string
	Usage() string
}

// option for subcommands
type option struct {
	name  string
	usage string
}

// Name return option name
func (o *option) Name() string {
	return o.name
}

// Usage return option usage
func (o *option) Usage() string {
	return o.usage
}

// StringOption wrap string option
type StringOption struct {
	option
	defaultVal string
}

// RegisterFlag register option as flag
func (o *StringOption) RegisterFlag() {
	flag.String(o.name, o.defaultVal, o.usage)
}

// UintOption  wrap uint option
type UintOption struct {
	option
	defaultVal uint
}

// RegisterFlag register option as flag
func (o *UintOption) RegisterFlag() {
	flag.Uint(o.name, o.defaultVal, o.usage)
}

// BoolOption  wrap bool option
type BoolOption struct {
	option
	defaultVal bool
}

// RegisterFlag register option as flag
func (o *BoolOption) RegisterFlag() {
	flag.Bool(o.name, o.defaultVal, o.usage)
}

// Float64Option  wrap float64 option
type Float64Option struct {
	option
	defaultVal float64
}

// RegisterFlag register option as flag
func (o *Float64Option) RegisterFlag() {
	flag.Float64(o.name, o.defaultVal, o.usage)
}
