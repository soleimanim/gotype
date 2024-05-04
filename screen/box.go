package screen

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/styles"
)

type TextAlignment string

const (
	TextAlignmentCenter TextAlignment = "center"
	TextAlignmentRight  TextAlignment = "right"
	TextAlignmentLeft   TextAlignment = "left"
)

type BoxTitle struct {
	Title     string
	Alignment TextAlignment
}

func DrawBox(position BufferPosition, size BufferSize, screen tcell.Screen, title BoxTitle, bgColor tcell.Color) {
	borderStyle := styles.BorderStyle.Background(bgColor)
	// corners
	screen.SetContent(position.X, position.Y, '╭', nil, borderStyle)
	screen.SetContent(position.X+size.Width-1, position.Y, '╮', nil, borderStyle)
	screen.SetContent(position.X, position.Y+size.Height-1, '╰', nil, borderStyle)
	screen.SetContent(position.X+size.Width-1, position.Y+size.Height-1, '╯', nil, borderStyle)

	// horizontal lines
	for i := position.X + 1; i < position.X+size.Width-1; i += 1 {
		screen.SetContent(i, position.Y, '─', nil, borderStyle)
		screen.SetContent(i, position.Y+size.Height-1, '─', nil, borderStyle)
	}

	// vertical lines
	for i := position.Y + 1; i < position.Y+size.Height-1; i += 1 {
		screen.SetContent(position.X, i, '│', nil, borderStyle)
		screen.SetContent(position.X+size.Width-1, i, '│', nil, borderStyle)
	}
	for i := position.X + 1; i < position.X+size.Width-1; i += 1 {
		for j := position.Y + 1; j < position.Y+size.Height-1; j += 1 {
			screen.SetContent(i, j, ' ', nil, borderStyle)
		}
	}

	if title.Title != "" {
		t := fmt.Sprintf(" %s ", title.Title)
		x := position.X + 4
		if title.Alignment == TextAlignmentCenter {
			x = position.X + size.Width/2 - len(t)/2
		} else if title.Alignment == TextAlignmentRight {
			x = position.X + size.Width - len(t) - 4
		}

		DrawText(screen, t, &x, &position.Y, borderStyle)
	}
}
