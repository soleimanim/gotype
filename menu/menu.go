package menu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/content"
)

type MenuItemAction func(*content.Content)

type MenuStyle tcell.Style

var (
	MenuStyleDefault MenuStyle = MenuStyle(tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorLightGreen))
)

type MenuItem struct {
	Name   string
	Key    tcell.Key
	Style  MenuStyle
	Action MenuItemAction
}
