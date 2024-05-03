package buffers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/screen"
)

const DIALOG_BUFFER_ID = 3
const DialogButtonsSpacing = 2
const DialogHorizontalPadding = 5
const DialogVerticalPadding = 2

type StyledText struct {
	Text  string
	Style tcell.Style
}

type DialogButton struct {
	Label  string
	Key    tcell.Key
	Style  tcell.Style
	Action func() bool
}

type DialogBuffer struct {
	screen tcell.Screen
	window *screen.Window

	Title       string
	Description []StyledText
	Buttons     []DialogButton

	TitleStyle       tcell.Style
	DescriptionStyle tcell.Style
}

func (b DialogBuffer) Draw() {
	screenWidth, screenHeight := b.screen.Size()

	titleLen := len(b.Title)
	descLen := 0
	for _, st := range b.Description {
		descLen += len(st.Text)
	}

	buttonsLen := 0
	for _, button := range b.Buttons {
		buttonsLen += len(button.Label) + DialogButtonsSpacing
	}

	width := titleLen
	if descLen > width {
		width = descLen
	} else if buttonsLen > width {
		width = buttonsLen
	}

	if width > screenWidth/3*2 {
		width = screenWidth / 3 * 2
	}

	containerWidth := width
	width += DialogHorizontalPadding * 2
	containerX := screenWidth/2 - width/2 + DialogHorizontalPadding

	height := DialogVerticalPadding * 2
	height += titleLen/containerWidth + 1
	height += descLen/containerWidth + 1
	height += buttonsLen/containerWidth + 1

	x := screenWidth/2 - width/2
	y := screenHeight/2 - height/2
	screen.DrawBox(screen.BufferPosition{
		X: x,
		Y: y,
	}, screen.BufferSize{
		Width:  width,
		Height: height,
	}, b.screen, screen.BoxTitle{
		Title:     b.Title,
		Alignment: screen.TextAlignmentLeft,
	}, tcell.ColorReset)

	if titleLen < containerWidth {
		// x := containerX + containerWidth/2 - titleLen/2 - 1
		// screen.DrawText(b.screen, " "+b.Title+" ", &x, &y, b.TitleStyle)
		y += DialogVerticalPadding
	} else {
		y += DialogVerticalPadding
		b.drawString(b.Title, containerX, &y, containerWidth, b.TitleStyle)
		y += 1
	}
	b.drawStyledTexts(b.Description, containerX, &y, containerWidth)

	y = screenHeight/2 - height/2 + height - 2
	x = containerX + containerWidth/2 - buttonsLen/2 + DialogButtonsSpacing
	for _, button := range b.Buttons {
		screen.DrawText(b.screen, button.Label, &x, &y, button.Style)
		x += DialogButtonsSpacing
	}
}

func (_ DialogBuffer) GetID() int {
	return DIALOG_BUFFER_ID
}
func (b *DialogBuffer) SetScreen(s tcell.Screen) {
	b.screen = s
}
func (b DialogBuffer) HandleKeyEvent(ev *tcell.EventKey) {
	for _, button := range b.Buttons {
		if ev.Key() == button.Key {
			if button.Action() {
				b.window.RemoveBuffer(b.GetID())
			}
		}
	}
}
func (b *DialogBuffer) SetWindow(w *screen.Window) {
	b.window = w
}

func (b DialogBuffer) drawString(s string, containerX int, y *int, containerWidth int, style tcell.Style) {
	l := len(s)
	if l < containerWidth {
		x := containerX + containerWidth/2 - l/2
		screen.DrawText(b.screen, s, &x, y, style)
		return
	}

	lines := l/containerWidth + 1
	remainingText := s
	for line := range lines {
		if line == lines-1 {
			// last line
			x := containerX + containerWidth/2 - len(remainingText)/2
			screen.DrawText(b.screen, remainingText, &x, y, style)
			return
		}
		text := remainingText[:containerWidth]
		x := containerX
		screen.DrawText(b.screen, text, &x, y, style)

		*y += 1
		remainingText = remainingText[containerWidth:]
	}
}

func (b DialogBuffer) drawStyledTexts(s []StyledText, containerX int, y *int, containerWidth int) {
	x := containerX
	for _, t := range s {
		text := t.Text
		if len(text)+x-containerX > containerWidth {
			x = containerX
			*y = 0
		}
		screen.DrawText(b.screen, text, &x, y, t.Style)
	}
}
