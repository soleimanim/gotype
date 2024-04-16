package screen

import "github.com/gdamore/tcell/v2"

func DrawText(screen tcell.Screen, s string, x *int, y *int, style tcell.Style) {
	for _, r := range s {
		*x++
		screen.SetContent(*x, *y, r, nil, style)
	}
}

func DrawParagraphInContainer(screen tcell.Screen, containerX *int, y *int, maxWidth int, text string, style tcell.Style, center bool) {
	linesCount := 1 + (len(text) / (maxWidth - 2))
	remainingText := text
	for line := range linesCount {
		if line == linesCount-1 {
			// last line
			x := *containerX + 1
			if center {
				x = *containerX + maxWidth/2 - len(remainingText)/2
			}
			DrawText(screen, remainingText, &x, y, style)
		} else {
			str := remainingText[:maxWidth-2]
			remainingText = remainingText[maxWidth-2:]
			x := *containerX + 1
			DrawText(screen, str, &x, y, style)
		}
		*y++
	}
}
