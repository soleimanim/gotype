package screen

import (
	"errors"
	"log"

	"github.com/gdamore/tcell/v2"
)

type Window struct {
	screen  tcell.Screen
	buffers []Buffer
}

func NewWindow() Window {
	return Window{}
}

// Initialize the window
//
// Parameters:
//   - mainBuffer: The first buffer that window will draw
//
// Returns:
//   - error: nil in case of success
func (w *Window) Init(mainBuffer Buffer) error {
	if len(w.buffers) > 0 {
		return errors.New("window is already initialized")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	err = screen.Init()
	if err != nil {
		return err
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))

	screen.Clear()
	w.screen = screen

	w.buffers = make([]Buffer, 0)
	w.AppendBuffer(mainBuffer)

	return nil
}

func (w Window) Draw() error {
	if len(w.buffers) == 0 {
		return errors.New("no buffer to draw")
	}

	w.screen.Clear()
	for _, b := range w.buffers {
		b.Draw()
	}
	w.screen.Show()

	return nil
}

// Adds new buffer to the window
//
// Parameters:
//   - b: buffer to add
func (w *Window) AppendBuffer(b Buffer) {
	b.SetScreen(w.screen)
	b.SetWindow(w)
	w.buffers = append(w.buffers, b)
}

// Removes buffer from window
// Parameters:
//   - b: the buffer to be removed
func (w *Window) RemoveBuffer(i int) {
	buffers := make([]Buffer, 0)
	for _, buffer := range w.buffers {
		if buffer.GetID() == i {
			continue
		} else {
			buffers = append(buffers, buffer)
		}
	}
	w.buffers = buffers
}

func (w *Window) ReplaceBuffer(id int, b Buffer) {
	w.RemoveBuffer(id)
	w.AppendBuffer(b)
}

// Handles tcell screen HandleEvents
//
// Returns:
//   - bool: terminate signal received
func (w Window) HandleEvents() bool {
	event := w.screen.PollEvent()
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlC {
			return true
		}
		for _, b := range w.buffers {
			b.HandleKeyEvent(ev)
		}
	}

	return false
}

func (w Window) GetBufferByID(id int) Buffer {
	for _, b := range w.buffers {
		if b.GetID() == id {
			return b
		}
	}

	return nil
}

// Close the window
func (w Window) Close() {
	w.screen.Fini()
}
