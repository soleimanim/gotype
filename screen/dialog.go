package screen

// import (
// 	"github.com/gdamore/tcell/v2"
// )
//
// var DialogBorderStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorGray)
// var DialogBackgroundStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorLightBlue)
// var DialogTitleStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorLightBlue).Foreground(tcell.ColorBlack)
// var DialogDescriptionStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorLightBlue).Foreground(tcell.ColorDarkGray)
// var DialogButtonStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorPink).Foreground(tcell.ColorBlack)
// var DialogInputStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorBlack)
// var DialogInputLabelStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorLightBlue).Foreground(tcell.ColorDarkGreen)
// var DialogCursorStyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
//
// type DialogActionFunc func(string) bool
//
// type DialogActionButton struct {
// 	Label  string
// 	Key    tcell.Key
// 	Action DialogActionFunc
// }
//
// type Dialog struct {
// 	Title       DialogText
// 	Description DialogText
// 	InputLabel  string
// 	InputValue  string
// 	Buttons     []DialogActionButton
//
// 	shouldShowInputField bool
// }
//
// type DialogTextAlignment int
//
// const (
// 	DialogTextAlignmentCenter DialogTextAlignment = 0
// )
//
// type DialogText struct {
// 	Text      string
// 	Alignment DialogTextAlignment
// 	Style     tcell.Style
// }
//
// func NewDialog(title string, description string, buttons []DialogActionButton) Dialog {
// 	return Dialog{
// 		Title: DialogText{
// 			Text:      title,
// 			Alignment: DialogTextAlignmentCenter,
// 			Style:     DialogTitleStyle,
// 		},
// 		Description: DialogText{
// 			Text:      description,
// 			Alignment: DialogTextAlignmentCenter,
// 			Style:     DialogDescriptionStyle,
// 		},
// 		Buttons: buttons,
// 	}
// }
//
// func (d *Dialog) SetInputFieldEnabled() *Dialog {
// 	d.shouldShowInputField = true
// 	return d
// }
//
// func (d *Dialog) SetInputLabel(l string) *Dialog {
// 	d.InputLabel = l
// 	return d
// }
//
// func (d *Dialog) HandleEvent(ev *tcell.EventKey) bool {
// 	for _, button := range d.Buttons {
// 		if button.Key == ev.Key() {
// 			return button.Action(d.InputValue)
// 		}
// 	}
// 	if d.shouldShowInputField {
// 		if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
// 			if len(d.InputValue) > 0 {
// 				d.InputValue = d.InputValue[:len(d.InputValue)-1]
// 			}
// 		} else {
// 			d.InputValue = d.InputValue + string(ev.Rune())
// 		}
// 	}
//
// 	return false
// }
//
// func (d Dialog) DrawDialogOnScreen(screen tcell.Screen) {
// 	screenWidth, screenHeight := screen.Size()
// 	width := screenWidth / 2
// 	containerWidth := width - 4
// 	x := screenWidth/2 - width/2
//
// 	// draw dialog box
// 	numberOfLines := 8
// 	numberOfLines += linesCount(d.Title.Text, containerWidth)
// 	numberOfLines += linesCount(d.Description.Text, containerWidth)
// 	if d.shouldShowInputField {
// 		numberOfLines += 2
// 	}
//
// 	y := screenHeight/2 - numberOfLines/2
// 	for i := range width {
// 		screen.SetContent(x+i, y, ' ', nil, DialogBorderStyle)
// 		screen.SetContent(x+i, y+numberOfLines-1, ' ', nil, DialogBorderStyle)
// 	}
// 	for i := range numberOfLines - 2 {
// 		bgY := y + i + 1
// 		screen.SetContent(x, bgY, ' ', nil, DialogBorderStyle)
// 		screen.SetContent(x+width-1, bgY, ' ', nil, DialogBorderStyle)
// 		for j := range width - 2 {
// 			screen.SetContent(x+j+1, bgY, ' ', nil, DialogBackgroundStyle)
// 		}
// 	}
//
// 	// draw title
// 	y += 2 // one line extra space on top
// 	txtX := x
// 	DrawParagraphInContainer(screen, &txtX, &y, containerWidth, d.Title.Text, d.Title.Style.Bold(true), true)
//
// 	y += 1
// 	txtX = x
// 	DrawParagraphInContainer(screen, &txtX, &y, containerWidth, d.Description.Text, d.Description.Style, true)
//
// 	if d.shouldShowInputField {
// 		y += 2
// 		inputX := x + 1
// 		DrawText(screen, d.InputLabel+" ", &inputX, &y, DialogInputLabelStyle)
// 		DrawText(screen, d.InputValue, &inputX, &y, DialogInputStyle)
// 		DrawText(screen, " ", &inputX, &y, DialogCursorStyle)
// 	}
//
// 	y += 2
// 	buttonsLen := 0
// 	for i, button := range d.Buttons {
// 		buttonsLen += len(button.Label)
// 		if i != len(d.Buttons)-1 {
// 			buttonsLen += 2
// 		}
// 	}
// 	buttonsX := x + containerWidth/2 - buttonsLen/2
// 	for _, button := range d.Buttons {
// 		DrawText(screen, " "+button.Label+" ", &buttonsX, &y, DialogButtonStyle)
// 		buttonsX += 2
// 	}
// }
//
// func linesCount(s string, width int) (n int) {
// 	if s != "" {
// 		if len(s) > width {
// 			n += len(s)/(width) + 1
// 		} else {
// 			n++
// 		}
// 	}
//
// 	return n
// }
