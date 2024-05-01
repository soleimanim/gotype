package buffers

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/screen"
)

const ACTION_LINE_BUFFER_ID = 4

type ActionsLineMenu interface {
	GetKey() tcell.Key
}

type ActionsLineMenuSwitch struct {
	Label    string
	Selected bool
	Key      tcell.Key
	OnChange func(bool, *ActionsLineBuffer)
}

func (s *ActionsLineMenuSwitch) GetKey() tcell.Key {
	return s.Key
}

type ActionsLineBuffer struct {
	screen tcell.Screen
	window *screen.Window

	menuItems []ActionsLineMenu
}

func NewActionsLineBuffer() *ActionsLineBuffer {
	b := &ActionsLineBuffer{}
	testModeMenu := make([]*ActionsLineMenuSwitch, 0)
	menu25Words := &ActionsLineMenuSwitch{
		Label:    " 25 Words ",
		Selected: true,
		Key:      tcell.KeyF1,
	}
	menu25Words.OnChange = func(b bool, alb *ActionsLineBuffer) {
		if b {
			alb.disableAllTestModeSwitches(testModeMenu)
			menu25Words.Selected = true
			alb.changeTypingTestMode(TestMode25Words)
		}
	}
	testModeMenu = append(testModeMenu, menu25Words)
	menu50Words := &ActionsLineMenuSwitch{
		Label: " 50 Words ",
		Key:   tcell.KeyF2,
	}
	menu50Words.OnChange = func(b bool, alb *ActionsLineBuffer) {
		if b {
			alb.disableAllTestModeSwitches(testModeMenu)
			menu50Words.Selected = true
			alb.changeTypingTestMode(TestMode50Words)
		}
	}
	testModeMenu = append(testModeMenu, menu50Words)
	menu75Words := &ActionsLineMenuSwitch{
		Label: " 75 Words ",
		Key:   tcell.KeyF3,
	}
	menu75Words.OnChange = func(b bool, alb *ActionsLineBuffer) {
		if b {
			alb.disableAllTestModeSwitches(testModeMenu)
			menu75Words.Selected = true
			alb.changeTypingTestMode(TestMode75Words)
		}
	}
	testModeMenu = append(testModeMenu, menu75Words)
	menu100Words := &ActionsLineMenuSwitch{
		Label: " 100 Words ",
		Key:   tcell.KeyF4,
	}
	menu100Words.OnChange = func(enabled bool, alb *ActionsLineBuffer) {
		if enabled {
			alb.disableAllTestModeSwitches(testModeMenu)
			menu100Words.Selected = true
			b.changeTypingTestMode(TestMode100Words)
		}
	}
	testModeMenu = append(testModeMenu, menu100Words)

	for _, m := range testModeMenu {
		b.menuItems = append(b.menuItems, m)
	}
	return b
}

func (b *ActionsLineBuffer) Draw() {
	screenWidth, screenHeight := b.screen.Size()
	y := screenHeight - 1
	for i := range screenWidth {
		b.screen.SetContent(i, y, ' ', nil, tcell.StyleDefault.Background(tcell.NewRGBColor(238, 234, 211)))
	}

	// draw menuItems
	x := 0
	bgStyle := tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorWhite)
	for i, m := range b.menuItems {
		switch menu := m.(type) {
		case *ActionsLineMenuSwitch:
			if menu.Selected {
				screen.DrawText(b.screen, " ✔ ", &x, &y, bgStyle.Foreground(tcell.ColorRed))
			}
			screen.DrawText(b.screen, menu.Label, &x, &y, bgStyle)
			screen.DrawText(b.screen, " "+tcell.KeyNames[menu.Key]+" ", &x, &y, bgStyle.Foreground(tcell.ColorGreen).Bold(true))
		}

		if i != len(b.menuItems)-1 {
			screen.DrawText(b.screen, " | ", &x, &y, bgStyle.Foreground(tcell.ColorDarkCyan))
		}
	}
}
func (b *ActionsLineBuffer) GetID() int {
	return ACTION_LINE_BUFFER_ID
}
func (b *ActionsLineBuffer) SetScreen(s tcell.Screen) {
	b.screen = s
}
func (b *ActionsLineBuffer) HandleKeyEvent(ev *tcell.EventKey) {
	for _, m := range b.menuItems {
		switch menu := m.(type) {
		case *ActionsLineMenuSwitch:
			if menu.Key == ev.Key() {
				menu.Selected = !menu.Selected
				menu.OnChange(menu.Selected, b)
			}
		}
	}
}
func (b *ActionsLineBuffer) SetWindow(w *screen.Window) {
	b.window = w
}

func (b *ActionsLineBuffer) changeTypingTestMode(mode TestMode) {
	buffer := NewTypingTestBuffer(mode)
	log.Println("b.window in buffer is", b.window, b.screen)
	b.window.ReplaceBuffer(TYPING_BUFFER_IDENTIFIER, &buffer)
}

func (b *ActionsLineBuffer) disableAllTestModeSwitches(menu []*ActionsLineMenuSwitch) {
	for _, m := range menu {
		m.Selected = false
	}
}
