package screen

import "github.com/gdamore/tcell/v2"

type BufferSize struct {
	Width  int
	Height int
}

type BufferPosition struct {
	X int
	Y int
}
type Buffer interface {
	Draw()
	GetID() int
	SetScreen(tcell.Screen)
	HandleKeyEvent(*tcell.EventKey)
	SetWindow(*Window)
}
