package buffers

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/db"
	"github.com/soleimanim/gotype/logger"
	"github.com/soleimanim/gotype/screen"
)

const RECENT_TESTS_BUFFER_ID = 5

type RecentTestsBuffer struct {
	Size        screen.BufferSize
	Position    screen.BufferPosition
	Repository  db.Repository[db.TypingTestModel]
	RecentTests []db.TypingTestModel

	MinSpeed    float32
	MaxSpeed    float32
	MinAccuracy float32
	MaxAccuracy float32

	window *screen.Window
	screen tcell.Screen
}

func NewRecentTestsBuffer(position screen.BufferPosition, size screen.BufferSize, repository db.Repository[db.TypingTestModel]) RecentTestsBuffer {
	return RecentTestsBuffer{
		Size:       size,
		Position:   position,
		Repository: repository,
	}
}

func (b RecentTestsBuffer) Draw() {
	screen.DrawBox(b.Position, b.Size, b.screen, screen.BoxTitle{
		Title:     "Recent Tests",
		Alignment: screen.TextAlignmentLeft,
	}, tcell.ColorReset)

	if len(b.RecentTests) == 0 {
		logger.Println("no recent tests, printing no recent test message")
		style := tcell.StyleDefault.Foreground(tcell.ColorLightGray)
		message := "No recent tests to show."
		x := b.Position.X + b.Size.Width/2 - len(message)/2
		y := b.Position.Y + 2
		screen.DrawText(b.screen, message, &x, &y, style)
		return
	}

	sampleFormatTime := "2024-30-01 01:01"
	timeLen := len(sampleFormatTime)
	maxSpeedLen := len("speed")
	maxAccLen := len("accuracy")

	for _, t := range b.RecentTests {
		speed := fmt.Sprintf("%.2f WPS", t.Speed)
		acc := fmt.Sprintf("%.2f%%", t.Accuracy)
		if len(speed) > maxSpeedLen {
			maxSpeedLen = len(speed)
		}
		if len(acc) > maxAccLen {
			maxAccLen = len(acc)
		}
	}

	b.Position.Y += 1
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlack)
	remainingSpace := (b.Size.Width - maxAccLen - maxSpeedLen - timeLen - 2) / 2

	timeX := b.Position.X + 1
	speedX := timeX + timeLen + remainingSpace
	accX := speedX + maxSpeedLen + remainingSpace
	headerStyle := style.Bold(true)
	screen.DrawText(b.screen, "Date", &timeX, &b.Position.Y, headerStyle)
	screen.DrawText(b.screen, "Speed", &speedX, &b.Position.Y, headerStyle)
	screen.DrawText(b.screen, "Accuracy", &accX, &b.Position.Y, headerStyle)
	b.Position.Y += 1

	for _, t := range b.RecentTests {
		timeX := b.Position.X + 1
		speedX := timeX + timeLen + remainingSpace
		accX := speedX + maxSpeedLen + remainingSpace

		time := b.formatTime(t.TestDate)
		speed := fmt.Sprintf("%.2f WPS", t.Speed)
		acc := fmt.Sprintf("%.2f%%", t.Accuracy)

		screen.DrawText(b.screen, time, &timeX, &b.Position.Y, style)
		screen.DrawText(b.screen, speed, &speedX, &b.Position.Y, b.getSpeedStyle(t))
		screen.DrawText(b.screen, acc, &accX, &b.Position.Y, b.getAccStyle(t))
		b.Position.Y += 1
	}
}
func (_ *RecentTestsBuffer) GetID() int {
	return RECENT_TESTS_BUFFER_ID
}
func (b *RecentTestsBuffer) SetScreen(s tcell.Screen) {
	b.screen = s
}
func (b *RecentTestsBuffer) HandleKeyEvent(_ *tcell.EventKey) {

}
func (b *RecentTestsBuffer) SetWindow(w *screen.Window) {
	b.window = w
}

func (b *RecentTestsBuffer) Update() {
	rowsCount := b.Size.Height - 3
	tests, err := b.Repository.GetAll(rowsCount, 0)
	if err != nil {
		return
	}
	b.RecentTests = tests
	for i, t := range tests {
		if i == 0 {
			b.MinAccuracy = t.Accuracy
			b.MaxAccuracy = t.Accuracy
			b.MinSpeed = t.Speed
			b.MaxSpeed = t.Speed
			continue
		}
		if t.Speed > b.MaxSpeed {
			b.MaxSpeed = t.Speed
		}
		if t.Speed < b.MinSpeed {
			b.MinSpeed = t.Speed
		}
		if t.Accuracy > b.MaxAccuracy {
			b.MaxAccuracy = t.Accuracy
		}
		if t.Accuracy < b.MinAccuracy {
			b.MinAccuracy = t.Accuracy
		}
	}
}

func (b RecentTestsBuffer) normalize(num, minNum, maxNum float32) float32 {
	n := ((num - minNum) / (maxNum - minNum))
	logger.Println("calculating normalozation", num, minNum, maxNum, n)
	return n
}

func (b RecentTestsBuffer) normalizedSpeed(num float32) float32 {
	return b.normalize(num, b.MinSpeed, b.MaxSpeed)
}

func (b RecentTestsBuffer) normalizedAccuracy(num float32) float32 {
	return b.normalize(num, b.MinAccuracy, b.MaxAccuracy)
}

func (b RecentTestsBuffer) getSpeedStyle(t db.TypingTestModel) tcell.Style {
	normalizedSpeed := b.normalizedSpeed(t.Speed)
	speedStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorReset)
	if normalizedSpeed == 1 {
		speedStyle = speedStyle.Foreground(tcell.ColorDarkBlue).Background(tcell.ColorLightSkyBlue)
	} else if normalizedSpeed == 0 {
		speedStyle = speedStyle.Foreground(tcell.ColorDarkRed).Background(tcell.ColorRed)
	} else if normalizedSpeed > 0.9 {
		speedStyle = speedStyle.Foreground(tcell.ColorDarkGreen)
	} else if normalizedSpeed > 0.8 {
		speedStyle = speedStyle.Foreground(tcell.ColorGreen)
	} else if normalizedSpeed > 0.7 {
		speedStyle = speedStyle.Foreground(tcell.ColorLightGreen)
	} else if normalizedSpeed > 0.6 {
		speedStyle = speedStyle.Foreground(tcell.ColorDarkOrange)
	} else if normalizedSpeed > 0.5 {
		speedStyle = speedStyle.Foreground(tcell.ColorOrange)
	} else {
		speedStyle = speedStyle.Foreground(tcell.ColorRed)
	}
	return speedStyle
}

func (b RecentTestsBuffer) getAccStyle(t db.TypingTestModel) tcell.Style {
	accStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorReset)
	normalizedAcc := b.normalizedAccuracy(t.Accuracy)
	if normalizedAcc == 1 {
		accStyle = accStyle.Foreground(tcell.ColorDarkBlue).Background(tcell.ColorLightSkyBlue)
	} else if normalizedAcc == 0 {
		accStyle = accStyle.Foreground(tcell.ColorDarkRed).Background(tcell.ColorRed)
	} else if normalizedAcc >= 0.7 {
		accStyle = accStyle.Foreground(tcell.ColorGreen)
	} else if normalizedAcc >= 0.6 {
		accStyle = accStyle.Foreground(tcell.ColorOrange)
	} else {
		accStyle = accStyle.Foreground(tcell.ColorRed)
	}
	return accStyle
}

func (b RecentTestsBuffer) formatTime(t time.Time) string {
	d := time.Since(t)
	quantity := 0
	unit := "hour"
	if d.Hours() < 1 {
		return "Recently"
	} else if d.Hours() < 24 {
		quantity = int(d.Hours())
		unit = "hour"
	} else {
		quantity = int(d.Hours() / 24)
		unit = "day"
	}

	pluralSign := ""
	if quantity > 1 {
		pluralSign = "s"
	}

	return fmt.Sprintf("%d %s%s ago", quantity, unit, pluralSign)
}
