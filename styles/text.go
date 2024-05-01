package styles

import "github.com/gdamore/tcell/v2"

var TextPlaceHolderStyle = tcell.StyleDefault.Foreground(tcell.ColorDimGray).Background(tcell.ColorReset)
var TextPrimaryStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorReset).Bold(true)
var TextErrorStyle = tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorReset).Underline(true)
var TextLogoStyle = tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorSkyblue).Bold(true)
var TextHighlightStyle1 = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorOrange).Bold(true)
var TextHighlightStyle2 = tcell.StyleDefault.Foreground(tcell.ColorLightYellow).Background(tcell.ColorSeaGreen).Bold(true)
