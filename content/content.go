package content

import (
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gdamore/tcell/v2"
)

type TextStyle tcell.Style

var TextStyleInfo1 TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorLightBlue).Foreground(tcell.ColorDarkBlue))
var TextStyleMain TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite))
var TextStyleError TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.ColorRed))
var TextStylePlaceholder TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkGray))
var TextStyleResult TextStyle = TextStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen))

type Content struct {
	Text      []rune
	InputText []rune

	StartTime     time.Time
	EndTime       time.Time
	MistakesCount uint
	Completed     bool
	FinalSpeed    float32
}

// Initialize a new content struct
//
// Returns:
//   - Content
func NewContent() *Content {
	return &Content{
		Text: []rune(randomdata.Paragraph()),
	}
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
	c.Text = []rune(randomdata.Paragraph())
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
