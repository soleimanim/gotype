package buffers

import (
	"fmt"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/db"
	"github.com/soleimanim/gotype/logger"
	"github.com/soleimanim/gotype/screen"
	"github.com/soleimanim/gotype/styles"
)

const TYPING_BUFFER_IDENTIFIER = 1

type TestMode int

const (
	TestMode25Words   TestMode = 0
	TestMode50Words   TestMode = 1
	TestMode75Words   TestMode = 2
	TestMode100Words  TestMode = 3
	TestMode15Seconds TestMode = 4
	TestMode30Seconds TestMode = 5
	TestMode45Seconds TestMode = 6
	TestMode60Seconds TestMode = 7
)

type TestBufferMenuItem struct {
	Label  string
	Key    tcell.Key
	Action func(*TypingTestBuffer)
}

var menuItems []TestBufferMenuItem = []TestBufferMenuItem{
	{
		Label: "New Test:",
		Key:   tcell.KeyF1,
		Action: func(b *TypingTestBuffer) {
			b.window.RemoveBuffer(TYPING_BUFFER_IDENTIFIER)
			buff := NewTypingTestBuffer(b.Position, b.Size, b.Mode, b.Repository)
			b.window.AppendBuffer(&buff)
		},
	},
	{
		Label: "Repeat:",
		Key:   tcell.KeyF2,
		Action: func(b *TypingTestBuffer) {
			b.reset()
		},
	},
}

type TypingEvent struct {
	Time time.Time
	Key  *tcell.EventKey
}

type TypingTestBuffer struct {
	Size     screen.BufferSize
	Position screen.BufferPosition

	Mode       TestMode
	Repository db.Repository[db.TypingTestModel]

	y        int
	screen   tcell.Screen
	window   *screen.Window
	input    string
	testText string

	recordedEvents []TypingEvent
	isReplaying    bool

	speed         float32
	accuracy      float32
	mistakesCount int

	isFinished bool
}

func GetTypingTestBufferPositionAndSize(s tcell.Screen) (screen.BufferPosition, screen.BufferSize) {
	w, h := s.Size()
	return screen.BufferPosition{
			X: w / 4,
			Y: 0,
		}, screen.BufferSize{
			Width:  w * 3 / 4,
			Height: h / 2,
		}
}

func NewTypingTestBuffer(position screen.BufferPosition, size screen.BufferSize, mode TestMode, repository db.Repository[db.TypingTestModel]) TypingTestBuffer {
	wordsCount := 25
	switch mode {
	case TestMode25Words:
		wordsCount = 25
	case TestMode50Words:
		wordsCount = 50
	case TestMode75Words:
		wordsCount = 75
	case TestMode100Words:
		wordsCount = 100
	}

	testText := randomdata.Noun()
	for range wordsCount - 1 {
		testText += " " + randomdata.Noun()
	}

	b := TypingTestBuffer{
		Size:     size,
		Position: position,

		Mode:       mode,
		Repository: repository,

		y:        5,
		screen:   nil,
		input:    "",
		testText: testText,
	}

	return b
}

func (b *TypingTestBuffer) Draw() {
	screen.DrawBox(b.Position, b.Size, b.screen, screen.BoxTitle{
		Title:     "Typing Test",
		Alignment: screen.TextAlignmentLeft,
	})
	b.y = b.Position.Y + 1
	startX := b.Position.X + 2
	for i, r := range b.testText {
		style := styles.TextPlaceHolderStyle
		if len(b.input) > i {
			if b.input[i] == b.testText[i] {
				style = styles.TextPrimaryStyle
			} else {
				style = styles.TextErrorStyle
			}
		}
		b.screen.SetContent(startX, b.y, r, nil, style)
		if i == len(b.input) {
			b.screen.ShowCursor(startX, b.y)
		}
		startX += 1
		if r == ' ' && len(b.testText) > i+2 {
			remainingText := b.testText[i+1:]
			for j, rr := range remainingText {
				if rr == ' ' {
					if startX+j+2 > b.Position.X+b.Size.Width {
						b.y += 1
						startX = b.Position.X + 2
					}
					break
				}
			}
		}
	}

	b.drawMenu()

	if b.isFinished || b.isReplaying {
		b.screen.HideCursor()
	}
	if b.isReplaying {
		y := b.y
		y += 2
		x := 0

		screen.DrawText(b.screen, "Replaying...", &x, &y, styles.ForegroundStyle(tcell.ColorLightGray))
		screen.DrawText(b.screen, " Press ", &x, &y, styles.ForegroundStyle(tcell.ColorLightGray))
		screen.DrawText(b.screen, " X ", &x, &y, styles.Style(tcell.ColorRed, tcell.ColorLightPink))
		screen.DrawText(b.screen, " to cancel replaying. ", &x, &y, styles.ForegroundStyle(tcell.ColorLightGray))
	}
}

func (b *TypingTestBuffer) GetID() int {
	return TYPING_BUFFER_IDENTIFIER
}
func (b *TypingTestBuffer) SetScreen(screen tcell.Screen) {
	b.screen = screen
}

func (b *TypingTestBuffer) SetWindow(w *screen.Window) {
	b.window = w
}

func (b *TypingTestBuffer) HandleKeyEvent(ev *tcell.EventKey) {
	for _, m := range menuItems {
		if m.Key == ev.Key() {
			m.Action(b)
			return
		}
	}

	if b.isReplaying {
		if ev.Rune() == 'x' || ev.Rune() == 'X' {
			b.isReplaying = false
			b.showTestResult()
		}
		return
	}
	if b.isFinished {
		return
	}

	b.recordEvent(ev)
	b.applyKeyEvent(ev)

	b.updateTestInfo()
	b.updateStatusLine()

	if len(b.input) >= len(b.testText) {
		b.isFinished = true
		b.onFinished()
		b.showTestResult()
		return
	}

}

func (b *TypingTestBuffer) Replay() {
	b.isReplaying = true
	b.input = ""
	b.y = 5
	b.window.Draw()

	go func() {
		for i, ev := range b.recordedEvents {
			if !b.isReplaying {
				break
			}
			if i != 0 {
				prevTime := b.recordedEvents[i-1].Time
				currentTime := b.recordedEvents[i].Time
				wait := currentTime.Sub(prevTime)
				time.Sleep(wait)
			}

			b.applyKeyEvent(ev.Key)
			b.window.Draw()
		}
		b.isReplaying = false
		b.showTestResult()
	}()
}

func (b *TypingTestBuffer) applyKeyEvent(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if len(b.input) > 0 {
			b.input = b.input[:len(b.input)-1]
		}
		return
	}
	if r := ev.Rune(); r != 0 {
		// character
		b.input = b.input + string(r)

		// check mistakes
		index := len(b.input) - 1
		if !b.isReplaying && b.input[index] != b.testText[index] {
			b.mistakesCount += 1
		}
	}
}

func (b *TypingTestBuffer) recordEvent(ev *tcell.EventKey) {
	b.recordedEvents = append(b.recordedEvents, TypingEvent{
		Time: time.Now(),
		Key:  ev,
	})
}

func (b TypingTestBuffer) countWords(s string) int {
	if len(s) == 0 {
		return 0
	}
	words := 1
	for _, r := range s {
		if r == ' ' {
			words += 1
		}
	}
	return words
}

func (b TypingTestBuffer) updateStatusLine() {
	buffer := b.window.GetBufferByID(STATUS_LINE_BUFFER_ID)
	statusBuffer, ok := buffer.(*StatusLineBuffer)
	if !ok {
		return
	}

	statusBuffer.Speed = b.speed
	statusBuffer.Accuracy = b.accuracy
}

func (b *TypingTestBuffer) updateTestInfo() {
	// calculate speed
	if len(b.recordedEvents) > 0 {
		words := float32(b.countWords(b.input))
		// all := (b.countWords(b.testText))
		startTime := b.recordedEvents[0].Time
		duration := time.Since(startTime).Seconds()
		b.speed = 60 * words / float32(duration)

		if b.mistakesCount > len(b.input) {
			b.accuracy = 0
		} else {
			b.accuracy = float32(len(b.testText)-b.mistakesCount) / float32(len(b.testText)) * 100.0
		}
	}
}

func (b *TypingTestBuffer) showTestResult() {
	dialog := DialogBuffer{
		Title: "Typing Test Result",
		Description: []StyledText{
			{
				Text:  "Speed: ",
				Style: styles.BorderStyle,
			},
			{
				Text:  fmt.Sprintf(" %.2f ", b.speed),
				Style: styles.TextHighlightStyle1,
			},
			{
				Text:  "\tAccuracy: ",
				Style: styles.BorderStyle,
			},
			{
				Text:  fmt.Sprintf(" %.2f ", b.accuracy),
				Style: styles.TextHighlightStyle2,
			},
		},
		Buttons: []DialogButton{
			{
				Label: " New Test ⏎ ",
				Key:   tcell.KeyEnter,
				Style: styles.Style(tcell.ColorSkyblue, tcell.ColorWhite),
				Action: func() bool {
					b.window.RemoveBuffer(TYPING_BUFFER_IDENTIFIER)
					newBuffer := NewTypingTestBuffer(b.Position, b.Size, b.Mode, b.Repository)
					b.window.AppendBuffer(&newBuffer)
					return true
				},
			},
			{
				Label: " Replay ^R ",
				Key:   tcell.KeyCtrlR,
				Style: styles.Style(tcell.ColorWhite, tcell.ColorSkyblue),
				Action: func() bool {
					b.Replay()
					return true
				},
			},
		},
		TitleStyle: styles.ForegroundStyle(tcell.ColorBlack).Bold(true),
	}

	b.window.AppendBuffer(&dialog)
}

func (b *TypingTestBuffer) onFinished() {
	logger.Println("TypingTestBuffer OnFinished called, saving result to database")
	model := db.TypingTestModel{
		Speed:         b.speed,
		Accuracy:      b.accuracy,
		WordsCount:    b.getWordsCount(b.Mode),
		MistakesCount: uint(b.mistakesCount),
	}
	err := b.Repository.Create(&model)
	if err != nil {
		logger.Println("Error saving typing test result to database", err)
		return
	}
	rb := b.window.GetBufferByID(RECENT_TESTS_BUFFER_ID)
	if rb != nil {
		buffer, ok := rb.(*RecentTestsBuffer)
		if ok {
			buffer.Update()
		}
	}
	sb := b.window.GetBufferByID(STATS_BUFFER_ID)
	if sb != nil {
		buffer, ok := sb.(*StatsBuffer)
		if ok {
			buffer.Update()
		}
	}
}

func (b TypingTestBuffer) getWordsCount(mode TestMode) uint {
	switch mode {
	case TestMode25Words:
		return 25
	case TestMode50Words:
		return 50
	case TestMode75Words:
		return 75
	case TestMode100Words:
		return 100
	}

	return 25
}

func (b TypingTestBuffer) drawMenu() {
	y := b.Size.Height - 3
	x := b.Position.X + 1
	w := b.Size.Width - 2

	lineX := x
	style := styles.BorderStyle.Foreground(tcell.ColorLightGray)
	for range w {
		b.screen.SetContent(lineX, y, '─', nil, style)
		lineX++
	}

	spacing := 5
	startX := x
	y += 1
	menuCount := len(menuItems)

	for i, m := range menuItems {
		screen.DrawText(b.screen, m.Label, &startX, &y, styles.ForegroundStyle(tcell.ColorRoyalBlue))
		startX++
		screen.DrawText(b.screen, tcell.KeyNames[m.Key], &startX, &y, styles.ForegroundStyle(tcell.ColorPurple))
		if i != menuCount-1 {
			b.screen.SetContent(startX+spacing/2, y, '|', nil, styles.BorderStyle)
		}
		startX += spacing
	}
}

func (b *TypingTestBuffer) reset() {
	b.input = ""
	b.mistakesCount = 0
	b.recordedEvents = []TypingEvent{}
	b.isReplaying = false
	b.speed = 0
	b.accuracy = 0
	b.isFinished = false
	b.speed = 0
	b.accuracy = 0
	b.isFinished = false
}
