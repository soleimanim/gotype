package content

import (
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gdamore/tcell/v2"
)

type TextStyle tcell.Style

var TextStyleInfo1 TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack))
var TextStyleInfo2 TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorLightBlue).Foreground(tcell.ColorDarkBlue))
var TextStyleInfo3 TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorPink).Foreground(tcell.ColorPurple))
var TextStyleMain TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite))
var TextStyleCursor TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorSeaGreen).Foreground(tcell.ColorNavajoWhite))
var TextStyleError TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite))
var TextStylePlaceholder TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkGray))
var TextStyleResult TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorSeaGreen))

type HighlightMode string

const (
	HighlightModeOnlyColor     HighlightMode = "HighLightModeOnlyColor"
	HighlightModeColorAndValue HighlightMode = "HighlightModeColorAndValue"
)

type TextMode int

const (
	TextModeParagraph TextMode = 0
	TextModeWords     TextMode = 1
)

type Content struct {
	Text      []rune
	InputText []rune

	TextMode TextMode

	StartTime          time.Time
	EndTime            time.Time
	MistakesCount      uint
	Completed          bool
	FinalSpeed         float32
	ErrorHighlightMode HighlightMode

	paragraphsCount int
}

// Initialize a new content struct
//
// Returns:
//   - Content
func NewContent() *Content {
	c := Content{
		TextMode:           TextModeParagraph,
		ErrorHighlightMode: HighlightModeOnlyColor,
		paragraphsCount:    1,
	}
	c.Text = []rune(generateText(c))
	return &c
}

// Add user input character
// this function will start timer when first character has been inputed
// also will count wrong input characters
//
// Parameters:
//   - r: the input rune
func (content *Content) AddInput(r rune) {
	if (content.StartTime == time.Time{}) {
		// Firts character typed, start timer
		content.StartTime = time.Now()
	}

	if len(content.InputText) < len(content.Text) {
		content.InputText = append(content.InputText, r)
		if content.Text[len(content.InputText)-1] != r {
			content.MistakesCount += 1
		}
	}
}

// Remove last input character
func (content *Content) RemoveLastInput() {
	if len(content.InputText) == 0 {
		return
	}
	content.InputText = content.InputText[:len(content.InputText)-1]
}

// Calculate typing speed
//
// Returns:
//   - speed
func (content Content) GetSpeed() float32 {
	if content.FinalSpeed != 0 {
		return content.FinalSpeed
	}

	if (content.StartTime == time.Time{}) {
		return 0
	}
	if len(content.InputText) < 1 {
		return 0
	}

	duration := time.Since(content.StartTime).Seconds()
	wordsCount := 1
	for _, r := range content.InputText[:len(content.InputText)-1] {
		if r == ' ' {
			wordsCount += 1
		}
	}

	return 60 * float32(wordsCount) / float32(duration)
}

// Calculate accurecy
func (content Content) GetAccuracy() float32 {
	if len(content.InputText) == 0 {
		return 0
	}

	correctCount := len(content.InputText) - int(content.MistakesCount)
	if correctCount < 0 {
		correctCount = 0
	}
	return float32(correctCount) / float32(len(content.InputText)) * 100
}

// Returns spent time from first input till now
//
// Returns:
//   - uint
func (content Content) GetSpentSeconds() uint {
	if (content.StartTime == time.Time{}) {
		return 0
	}

	if content.IsCompleted() {
		return uint(content.EndTime.Sub(content.StartTime).Seconds())
	}

	return uint(time.Since(content.StartTime).Seconds())
}

// Reset content data
func (c *Content) Reset() {
	c.Text = []rune(generateText(*c))
	c.InputText = []rune{}
	c.StartTime = time.Time{}
	c.MistakesCount = 0
	c.EndTime = time.Time{}
	c.FinalSpeed = 0
	c.Completed = false
}

// Detect if typing text is completed
//
// Returns:
//   - bool: typing the text is completed or not
func (c *Content) IsCompleted() bool {
	if c.Completed {
		return true
	}
	result := len(c.Text) == len(c.InputText)
	if result {
		c.FinalSpeed = c.GetSpeed()
		c.EndTime = time.Now()
	}
	c.Completed = result
	return result
}

func (c *Content) ToggleErrorHighlightingMode() {
	if c.ErrorHighlightMode == HighlightModeColorAndValue {
		c.ErrorHighlightMode = HighlightModeOnlyColor
	} else {
		c.ErrorHighlightMode = HighlightModeColorAndValue
	}
}

func (cnt *Content) SetParagraphsCount(c int) *Content {
	cnt.paragraphsCount = c
	cnt.Reset()
	return cnt
}
func generateText(c Content) string {
	if c.TextMode == TextModeParagraph {
		str := ""
		for i := range c.paragraphsCount {
			str += randomdata.Paragraph()
			if i != c.paragraphsCount-1 {
				str += " "
			}
		}

		return str
	}

	count := 18
	s := ""
	for i := range count {
		w := randomdata.Noun()
		s += w
		if i != count-1 {
			s += " "
		}
	}
	return s
}
