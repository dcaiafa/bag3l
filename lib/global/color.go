package global

import (
	"fmt"

	"github.com/dcaiafa/bag3l/internal/vm"
	fatihcolor "github.com/fatih/color"
)

// ColorMod is a print modifier that applies color formatting.
type ColorMod struct {
	Color *fatihcolor.Color
}

func (m *ColorMod) String() string    { return "<color>" }
func (m *ColorMod) Type() string      { return "color" }
func (m *ColorMod) Traits() vm.Traits { return vm.TraitNone }

func color0(m *vm.VM, args []string) (*ColorMod, error) {
	attribs := make([]fatihcolor.Attribute, 0, len(args))

	for _, arg := range args {
		attrib, err := ColorAttribute(arg)
		if err != nil {
			return nil, err
		}

		attribs = append(attribs, attrib)
	}

	c := &ColorMod{
		Color: fatihcolor.New(attribs...),
	}

	return c, nil
}

// ColorAttribute converts a color name string to a fatih/color attribute.
func ColorAttribute(v string) (attrib fatihcolor.Attribute, err error) {
	switch v {
	case "reset":
		attrib = fatihcolor.Reset
	case "bold":
		attrib = fatihcolor.Bold
	case "faint":
		attrib = fatihcolor.Faint
	case "italic":
		attrib = fatihcolor.Italic
	case "underline":
		attrib = fatihcolor.Underline
	case "blinkslow":
		attrib = fatihcolor.BlinkSlow
	case "blinkrapid":
		attrib = fatihcolor.BlinkRapid
	case "reversevideo":
		attrib = fatihcolor.ReverseVideo
	case "concealed":
		attrib = fatihcolor.Concealed
	case "crossedout":
		attrib = fatihcolor.CrossedOut
	case "black":
		attrib = fatihcolor.FgBlack
	case "red":
		attrib = fatihcolor.FgRed
	case "green":
		attrib = fatihcolor.FgGreen
	case "yellow":
		attrib = fatihcolor.FgYellow
	case "blue":
		attrib = fatihcolor.FgBlue
	case "magenta":
		attrib = fatihcolor.FgMagenta
	case "cyan":
		attrib = fatihcolor.FgCyan
	case "white":
		attrib = fatihcolor.FgWhite
	case "hiblack":
		attrib = fatihcolor.FgHiBlack
	case "hired":
		attrib = fatihcolor.FgHiRed
	case "higreen":
		attrib = fatihcolor.FgHiGreen
	case "hiyellow":
		attrib = fatihcolor.FgHiYellow
	case "hiblue":
		attrib = fatihcolor.FgHiBlue
	case "himagenta":
		attrib = fatihcolor.FgHiMagenta
	case "hicyan":
		attrib = fatihcolor.FgHiCyan
	case "hiwhite":
		attrib = fatihcolor.FgHiWhite
	case "bgblack":
		attrib = fatihcolor.BgBlack
	case "bgred":
		attrib = fatihcolor.BgRed
	case "bggreen":
		attrib = fatihcolor.BgGreen
	case "bgyellow":
		attrib = fatihcolor.BgYellow
	case "bgblue":
		attrib = fatihcolor.BgBlue
	case "bgmagenta":
		attrib = fatihcolor.BgMagenta
	case "bgcyan":
		attrib = fatihcolor.BgCyan
	case "bgwhite":
		attrib = fatihcolor.BgWhite
	case "bghiblack":
		attrib = fatihcolor.BgHiBlack
	case "bghired":
		attrib = fatihcolor.BgHiRed
	case "bghigreen":
		attrib = fatihcolor.BgHiGreen
	case "bghiyellow":
		attrib = fatihcolor.BgHiYellow
	case "bghiblue":
		attrib = fatihcolor.BgHiBlue
	case "bghimagenta":
		attrib = fatihcolor.BgHiMagenta
	case "bghicyan":
		attrib = fatihcolor.BgHiCyan
	case "bghiwhite":
		attrib = fatihcolor.BgHiWhite
	default:
		return fatihcolor.Reset, fmt.Errorf("invalid color attribute %v", v)
	}
	return attrib, nil
}
