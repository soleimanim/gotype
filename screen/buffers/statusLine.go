package buffers

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/screen"
	"github.com/soleimanim/gotype/styles"
)

const STATUS_LINE_BUFFER_ID = 2

type StatusLineBuffer struct {
	window *screen.Window
	screen tcell.Screen
	y      int

	Speed    float32
	Accuracy float32
}

func NewStatusLineBuffer() StatusLineBuffer {
	return StatusLineBuffer{}
}

func (b StatusLineBuffer) Draw() {
	screenWidth, _ := b.screen.Size()
	logoText := " GoType "
	startX := screenWidth/2 - len(logoText)/2
	screen.DrawText(b.screen, logoText, &startX, &b.y, styles.TextLogoStyle)

	startX = 0
	speedText := fmt.Sprintf(" Speed: %.2f ", b.Speed)
	screen.DrawText(b.screen, speedText, &startX, &b.y, styles.TextHighlightStyle1)

	accText := fmt.Sprintf(" Accuracy: %.2f ", b.Accuracy)
	startX = screenWidth - len(accText)
	screen.DrawText(b.screen, accText, &startX, &b.y, styles.TextHighlightStyle2)
}
func (b StatusLineBuffer) GetID() int {
	return STATUS_LINE_BUFFER_ID
}
func (b *StatusLineBuffer) SetScreen(s tcell.Screen) {
	b.screen = s
}
func (statuslinebuffer StatusLineBuffer) HandleKeyEvent(_ *tcell.EventKey) {
	return
}
func (b *StatusLineBuffer) SetWindow(w *screen.Window) {
	b.window = w
}
