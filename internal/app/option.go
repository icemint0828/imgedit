package app

import (
	"flag"
)

var OptionVertical = &BoolOption{
	option: option{
		name:  "vertical",
		usage: "specify direction as vertical. default horizon.",
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
		usage: "ratio. if ratio is set, width and height are ignored.",
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
		usage: "x length.",
	},
	defaultVal: 0,
}
var OptionY = &UintOption{
	option: option{
		name:  "y",
		usage: "y length.",
	},
	defaultVal: 0,
}
var OptionText = &StringOption{
	option: option{
		name:  "text",
		usage: "text.",
	},
	defaultVal: "",
}
var OptionTtf = &StringOption{
	option: option{
		name:  "ttf",
		usage: "ttf file path.",
	},
	defaultVal: "",
}
var OptionSize = &UintOption{
	option: option{
		name:  "size",
		usage: "font size.",
	},
	defaultVal: 0,
}
var OptionColor = &StringOption{
	option: option{
		name:  "color",
		usage: "font color with string (back, white, red, blue, green). or specify by color code(like #FF0000)",
	},
	defaultVal: "",
}
var OptionMode = &StringOption{
	option: option{
		name:  "mode",
		usage: "filter color(sepia, gray).",
	},
	defaultVal: "",
}

// Option for subcommands
type Option interface {
	RegisterFlag()
	Name() string
	Usage() string
	IsSet() bool
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

// IsSet return flag is set
func (o *option) IsSet() bool {
	isSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == o.name {
			isSet = true
			return
		}
	})
	return isSet
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

// String return flag value as string
func (o *StringOption) String() string {
	return flag.Lookup(o.Name()).Value.(flag.Getter).Get().(string)
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

// Uint return flag value as uint
func (o *UintOption) Uint() uint {
	return flag.Lookup(o.Name()).Value.(flag.Getter).Get().(uint)
}

// Int return flag value as int
func (o *UintOption) Int() int {
	return int(o.Uint())
}

// Float64 return flag value as float64
func (o *UintOption) Float64() float64 {
	return float64(o.Uint())
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

// Bool return flag value as bool
func (o *BoolOption) Bool() bool {
	return flag.Lookup(o.Name()).Value.(flag.Getter).Get().(bool)
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

// Float64 return flag value as float64
func (o *Float64Option) Float64() float64 {
	return flag.Lookup(o.Name()).Value.(flag.Getter).Get().(float64)
}
