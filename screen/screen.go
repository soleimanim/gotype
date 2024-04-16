package screen

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/content"
	"github.com/soleimanim/gotype/menu"
)

type ScreenStyle tcell.Style

var ScreenStyleDefault ScreenStyle = ScreenStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))

type Screen struct {
	screen     tcell.Screen
	menuItems  []menu.MenuItem
	showCursor bool

	activeDialog Dialog
	showDialog   bool
}

// Initialize a new screen
//
// Parameters:
//   - style: style to be used for screen
//
// Returns:
//   - Screen: the screen struct
//   - error: nil if initializing screen is successful
func NewScreen(style ScreenStyle, menuItems []menu.MenuItem) (Screen, error) {
	screen := Screen{
		menuItems:  menuItems,
		showCursor: true,
	}
	s, err := tcell.NewScreen()
	if err != nil {
		return screen, err
	}

	err = s.Init()
	if err != nil {
		return screen, err
	}

	s.SetStyle(tcell.Style(style))
	screen.screen = s

	return screen, nil
}

func (s *Screen) AddMenuItem(m menu.MenuItem) {
	s.menuItems = append(s.menuItems, m)
}

func (s *Screen) ToggleCursor() {
	s.showCursor = !s.showCursor
}

// Draw screen
//
// Parameters:
//   - content: the content struct to be drawn on the screen
//   - menuItems: List of menu items
func (s Screen) Draw(c *content.Content) {
	s.screen.Clear()

	y := 1
	drawInfo(s.screen, &y, c)

	y += 2
	x := 0

	screenWidth, _ := s.screen.Size()
	for index, r := range c.Text {
		// wrap text to next line
		if r == ' ' && index < len(c.Text)-1 {
			// check if we have enough room for next world
			for i, n := range c.Text[index+1:] {
				if n == ' ' {
					if x+i > screenWidth-2 {
						y += 1
						x = 0
					}
					break
				}
			}
		}
		if len(c.InputText) >= index+1 {
			input := c.InputText[index]
			if input == r {
				s.screen.SetContent(x, y, input, nil, tcell.Style(content.TextStyleMain))
			} else {
				ch := input
				if c.ErrorHighlightMode == content.HighlightModeOnlyColor {
					ch = r
				}
				s.screen.SetContent(x, y, ch, nil, tcell.Style(content.TextStyleError))
			}
		} else if index == len(c.InputText) && s.showCursor {
			s.screen.SetContent(x, y, r, nil, tcell.Style(content.TextStyleCursor))
			// s.screen.SetCursorStyle(tcell.CursorStyleDefault)
			// Handle cursor by setting backround color because tcell does not support cursor color yet
			// s.screen.ShowCursor(x, y)
		} else {
			// if !s.showCursor {
			// s.screen.HideCursor()
			// }
			s.screen.SetContent(x, y, r, nil, tcell.Style(content.TextStylePlaceholder))
		}

		x++
	}

	if c.IsCompleted() {
		y += 2
		text := fmt.Sprintf("Your typing speed is %.2f WPS Accuracy: %.2f%%", c.GetSpeed(), c.GetAccuracy())
		resultX := 0
		DrawText(s.screen, text, &resultX, &y, tcell.Style(content.TextStyleResult))
	}
	drawMenu(s)

	if s.showDialog {
		s.activeDialog.DrawDialogOnScreen(s.screen)
	}

	s.screen.Show()
}

// Read input events
func (s Screen) ReadEvent() tcell.Event {
	return s.screen.PollEvent()
}

func (s *Screen) HandleEvent(event *tcell.EventKey, c *content.Content) {
	if s.showDialog {
		r := s.activeDialog.HandleEvent(event)
		if r {
			s.showDialog = false
			s.activeDialog = Dialog{}
		}
		return
	}

	if (event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2) && !c.IsCompleted() {
		c.RemoveLastInput()
		return
	}

	for _, m := range s.menuItems {
		if m.Key == event.Key() {
			m.Action(c)
			s.Draw(c)
			return
		}
	}

	if !c.IsCompleted() {
		input := event.Rune()
		c.AddInput(input)
	}
}

func (s *Screen) SetDialog(d Dialog) {
	s.activeDialog = d
	s.showDialog = true
}

func (s Screen) Fini() {
	s.screen.Fini()
}
func drawInfo(screen tcell.Screen, y *int, c *content.Content) {
	x := 0

	timeSpentText := fmt.Sprintf(" Time: %d Seconds", c.GetSpentSeconds())
	DrawText(screen, timeSpentText, &x, y, tcell.Style(content.TextStyleInfo1))

	logoText := " GOTYPE "
	screenWith, _ := screen.Size()
	logoX := screenWith/2 - len(logoText)/2
	DrawText(screen, logoText, &logoX, y, tcell.Style(content.TextStyleInfo2))

	speedText := fmt.Sprintf(" Speed: %.2f WPS ", c.GetSpeed())
	speedX := screenWith - len(speedText) - 2
	DrawText(screen, speedText, &speedX, y, tcell.Style(content.TextStyleInfo3))
}

func drawMenu(s Screen) {
	_, screenHeight := s.screen.Size()
	y := screenHeight - 1
	x := 0

	for _, m := range s.menuItems {
		DrawText(s.screen, fmt.Sprintf(" %s ", m.Name), &x, &y, tcell.Style(m.Style))
		x += 2
	}
}
