package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/content"
	"github.com/soleimanim/gotype/menu"
	"github.com/soleimanim/gotype/screen"
)

func main() {
	s, err := screen.NewScreen(screen.ScreenStyleDefault, []menu.MenuItem{
		{
			Name:  "Retry (Ctrl + R)",
			Key:   tcell.KeyCtrlR,
			Style: menu.MenuStyleDefault,
			Action: func(c *content.Content) {
				c.Reset()
			},
		},
		{
			Name:  "Toggle Error Highlighting Mode (Ctrl + T)",
			Key:   tcell.KeyCtrlT,
			Style: menu.MenuStyleDefault,
			Action: func(c *content.Content) {
				c.ToggleErrorHighlightingMode()
			},
		},
	})
	s.AddMenuItem(
		menu.MenuItem{
			Name:  "Toggle Cursor (Ctrl + O)",
			Key:   tcell.KeyCtrlO,
			Style: menu.MenuStyleDefault,
			Action: func(c *content.Content) {
				s.ToggleCursor()
			},
		},
	)
	if err != nil {
		panic(err)
	}

	c := content.NewContent()

	for {
		s.Draw(c)
		event := s.ReadEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else {
				s.HandleEvent(ev, c)
			}
		}
	}
}
