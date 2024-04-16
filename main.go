package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/content"
	"github.com/soleimanim/gotype/menu"
	"github.com/soleimanim/gotype/screen"
)

func main() {
	s, err := screen.NewScreen(screen.ScreenStyleDefault, []menu.MenuItem{
		{
			Name:  "Retry ^R",
			Key:   tcell.KeyCtrlR,
			Style: menu.MenuStyleDefault,
			Action: func(c *content.Content) {
				c.Reset()
			},
		},
		// {
		// 	Name:  "Toggle Error Highlighting Mode ^T",
		// 	Key:   tcell.KeyCtrlT,
		// 	Style: menu.MenuStyleDefault,
		// 	Action: func(c *content.Content) {
		// 		c.ToggleErrorHighlightingMode()
		// 	},
		// },
	})
	s.AddMenuItem(
		menu.MenuItem{
			Name:  "Toggle Cursor ^O",
			Key:   tcell.KeyCtrlO,
			Style: menu.MenuStyleDefault,
			Action: func(c *content.Content) {
				s.ToggleCursor()
			},
		},
	)
	s.AddMenuItem(menu.MenuItem{
		Name:  "Set Paragraphs Count ^P",
		Key:   tcell.KeyCtrlP,
		Style: menu.MenuStyleDefault,
		Action: func(c *content.Content) {
			dialog := screen.NewDialog("Paragraphs Count", "Set how many paragraphs of text must be generated.", []screen.DialogActionButton{
				{
					Label: "OK ⏎",
					Key:   tcell.KeyEnter,
					Action: func(value string) bool {
						n, err := strconv.ParseInt(value, 10, 0)
						if err != nil {
							return false
						}
						c.SetParagraphsCount(int(n))
						return true
					},
				},
			})
			dialog.SetInputFieldEnabled()
			dialog.SetInputLabel("Paragraphs:")
			s.SetDialog(dialog)
		},
	})
	if err != nil {
		panic(err)
	}

	defer s.Fini()

	c := content.NewContent()

	eventCh := make(chan tcell.Event, 0)
	go readInput(eventCh, s)

	for {
		select {
		case ev := <-eventCh:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else {
					s.HandleEvent(ev, c)
				}
			}
			s.Draw(c)
			go readInput(eventCh, s)
		default:
			s.Draw(c)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func readInput(ch chan tcell.Event, s screen.Screen) {
	event := s.ReadEvent()
	ch <- event
}
