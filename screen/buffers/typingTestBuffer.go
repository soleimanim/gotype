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

type TypingEvent struct {
	Time time.Time
	Key  *tcell.EventKey
}

type TypingTestBuffer struct {
	Mode        TestMode
	Repository  db.Repository[db.TypingTestModel]
	recentTests []db.TypingTestModel

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

func NewTypingTestBuffer(mode TestMode, repository db.Repository[db.TypingTestModel]) TypingTestBuffer {
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
		Mode:       mode,
		Repository: repository,

		y:        5,
		screen:   nil,
		input:    "",
		testText: testText,
	}
	b.updateRecentTests()

	return b
}

func (b *TypingTestBuffer) Draw() {
	screenWidth, _ := b.screen.Size()
	b.y = 5
	startX := 0
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
					if startX+j+1 > screenWidth {
						b.y += 1
						startX = 0
					}
					break
				}
			}
		}
	}

	if b.isFinished || b.isReplaying {
		b.screen.HideCursor()
	}
	if b.isReplaying {
		y := b.y
		y += 2
		x := 0

		screen.DrawText(b.screen, "Replaying...", &x, &y, tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorLightGrey))
		screen.DrawText(b.screen, " Press ", &x, &y, tcell.StyleDefault.Foreground(tcell.ColorLightGray))
		screen.DrawText(b.screen, " X ", &x, &y, tcell.StyleDefault.Background(tcell.ColorLightPink).Foreground(tcell.ColorRed))
		screen.DrawText(b.screen, " to cancel replaying. ", &x, &y, tcell.StyleDefault.Foreground(tcell.ColorLightGray))
	}

	b.y += 3
	for _, test := range b.recentTests {
		x := 5
		screen.DrawText(b.screen, fmt.Sprintf("Speed: %.2f  Acc: %.2f  Words Count: %d  Mistakes Count %d", test.Speed, test.Accuracy, test.WordsCount, test.MistakesCount), &x, &b.y, tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorReset))
		b.y += 1
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
				Style: tcell.StyleDefault.Background(tcell.ColorWhite).Background(tcell.ColorSkyblue),
				Action: func() bool {
					b.window.RemoveBuffer(TYPING_BUFFER_IDENTIFIER)
					newBuffer := NewTypingTestBuffer(b.Mode, b.Repository)
					b.window.AppendBuffer(&newBuffer)
					return true
				},
			},
			{
				Label: " Replay ^R ",
				Key:   tcell.KeyCtrlR,
				Style: tcell.StyleDefault.Background(tcell.ColorWhite).Background(tcell.ColorSkyblue),
				Action: func() bool {
					b.Replay()
					return true
				},
			},
		},
		TitleStyle: tcell.StyleDefault.Bold(true).Foreground(tcell.ColorBlack),
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
	b.updateRecentTests()
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

func (b *TypingTestBuffer) updateRecentTests() {
	models, err := b.Repository.GetAll(10, 0)
	logger.Println("Read recent typing tests from database", len(models), err)
	if err == nil {
		b.recentTests = models
	}
}
