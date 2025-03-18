package color

import (
	"strconv"

	"github.com/fatih/color"
)

var (
	colorONE   = color.RGB(0xFF, 0x00, 0x00).Add(color.Bold).Sprint("1")
	colorTWO   = color.RGB(0x00, 0xFF, 0x00).Add(color.Bold).Sprint("2")
	colorTHREE = color.RGB(0x00, 0x00, 0xFF).Add(color.Bold).Sprint("3")
	colorFOUR  = color.RGB(0xFF, 0xFF, 0x00).Add(color.Bold).Sprint("4")
	colorFIVE  = color.RGB(0xFF, 0x00, 0xFF).Add(color.Bold).Sprint("5")
	colorSIX   = color.RGB(0x00, 0xFF, 0xFF).Add(color.Bold).Sprint("6")
)

func Colorize(v uint8) string {
	switch v {
	case 1:
		return colorONE
	case 2:
		return colorTWO
	case 3:
		return colorTHREE
	case 4:
		return colorFOUR
	case 5:
		return colorFIVE
	case 6:
		return colorSIX
	default:
		return unknown(strconv.Itoa(int(v)))
	}
}

func unknown(str string) string {
	return color.BgRGB(0xFF, 0x00, 0x00).Add(color.Bold).Sprint(str)
}
