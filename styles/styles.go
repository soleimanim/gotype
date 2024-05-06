package styles

import "github.com/gdamore/tcell/v2"

func Style(fg tcell.Color, bg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}

func ForegroundStyle(c tcell.Color) tcell.Style {
	return Style(c, tcell.ColorReset)
}

func StyleReset() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
}
