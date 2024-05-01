package screen

import "github.com/gdamore/tcell/v2"

type Buffer interface {
	Draw()
	GetID() int
	SetScreen(tcell.Screen)
	HandleKeyEvent(*tcell.EventKey)
	SetWindow(*Window)
}
