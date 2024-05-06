package styles

import "github.com/gdamore/tcell/v2"

var TextPlaceHolderStyle = Style(tcell.ColorReset, tcell.ColorReset)
var TextPrimaryStyle = Style(tcell.ColorSlateGray, tcell.ColorReset).Bold(true)
var TextErrorStyle = Style(tcell.ColorRed, tcell.ColorReset).Underline(true)
var TextLogoStyle = Style(tcell.ColorBlue, tcell.ColorSkyblue).Bold(true)
var TextHighlightStyle1 = Style(tcell.ColorWhite, tcell.ColorOrange).Bold(true)
var TextHighlightStyle2 = Style(tcell.ColorLightYellow, tcell.ColorSeaGreen).Bold(true)
