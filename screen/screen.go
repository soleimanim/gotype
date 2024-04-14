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
	screen    tcell.Screen
	menuItems []menu.MenuItem
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
		menuItems: menuItems,
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
				s.screen.SetContent(x, y, input, nil, tcell.Style(content.TextStyleError))
			}
		} else {
			s.screen.SetContent(x, y, r, nil, tcell.Style(content.TextStylePlaceholder))
		}

		x++
	}

	if c.IsCompleted() {
		y += 2
		text := fmt.Sprintf("Your typing speed is %.2f WPS Accuracy: %.2f%%", c.GetSpeed(), c.GetAccuracy())
		resultX := 0
		drawText(s.screen, &resultX, &y, text, content.TextStyleResult)
	}
	drawMenu(s)

	s.screen.Show()
}

// Read input events
func (s Screen) ReadEvent() tcell.Event {
	return s.screen.PollEvent()
}

func (s Screen) HandleEvent(event *tcell.EventKey, c *content.Content) {
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

func drawInfo(screen tcell.Screen, y *int, c *content.Content) {
	x := 0

	timeSpentText := fmt.Sprintf(" Time: %d Seconds", c.GetSpentSeconds())
	drawText(screen, &x, y, timeSpentText, content.TextStyleInfo1)

	logoText := " GOTYPE "
	screenWith, _ := screen.Size()
	logoX := screenWith/2 - len(logoText)/2
	drawText(screen, &logoX, y, logoText, content.TextStyleInfo1)

	speedText := fmt.Sprintf(" Speed: %.2f WPS ", c.GetSpeed())
	speedX := screenWith - len(speedText) - 2
	drawText(screen, &speedX, y, speedText, content.TextStyleInfo1)
}

func drawMenu(s Screen) {
	_, screenHeight := s.screen.Size()
	y := screenHeight - 1
	x := 0

	for _, m := range s.menuItems {
		drawText(s.screen, &x, &y, m.Name, content.TextStyle(m.Style))
	}
}

func drawText(screen tcell.Screen, x *int, y *int, text string, style content.TextStyle) {
	screenWidth, _ := screen.Size()
	if *x+len(text) > screenWidth-1 {
		*y += 1
		*x = 0
	}

	for index, r := range text {
		screen.SetContent(*x+index, *y, r, nil, tcell.Style(style))
	}
}
