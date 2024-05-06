package buffers

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/soleimanim/gotype/db"
	"github.com/soleimanim/gotype/screen"
	"github.com/soleimanim/gotype/styles"
)

const STATS_BUFFER_ID = 6

type Stat struct {
	Key   string
	Value any
}

func (s Stat) FormatValue(format string) string {
	switch v := s.Value.(type) {
	case float32, int:
		if v == -1 {
			return "N/A"
		}
		return fmt.Sprintf(format, s.Value)
	}

	return fmt.Sprintf(format, s.Value)
}

type StatsBuffer struct {
	window     *screen.Window
	screen     tcell.Screen
	Repository db.Repository[db.TypingTestModel]

	Position screen.BufferPosition
	Size     screen.BufferSize

	TestsCount         Stat
	SecondsSpentTyping Stat
	BestSpeed          Stat
	BestAccuracy       Stat
	WordsCount         Stat
	AvgSpeed           Stat
	AvgAccuracy        Stat
}

func NewStatsBuffer(position screen.BufferPosition, size screen.BufferSize, repository db.Repository[db.TypingTestModel]) StatsBuffer {
	return StatsBuffer{
		Repository: repository,
		Position:   position,
		Size:       size,
	}
}

func (b StatsBuffer) Draw() {
	screen.DrawBox(b.Position, b.Size, b.screen, screen.BoxTitle{
		Title:     "Stats",
		Alignment: screen.TextAlignmentLeft,
	})

	keyStyle := styles.ForegroundStyle(tcell.ColorReset)
	valStyle := styles.ForegroundStyle(tcell.ColorDarkGray).Bold(true)

	val := b.BestSpeed.FormatValue("%.2f WPS")
	b.drawKeyVal(b.BestSpeed.Key, val, keyStyle, valStyle)

	val = b.BestAccuracy.FormatValue("%.2f %%")
	b.drawKeyVal(b.BestAccuracy.Key, val, keyStyle, valStyle)

	val = b.TestsCount.FormatValue("%d Tests")
	b.drawKeyVal(b.TestsCount.Key, val, keyStyle, valStyle)

	val = b.WordsCount.FormatValue("%d Words")
	b.drawKeyVal(b.WordsCount.Key, val, keyStyle, valStyle)

	val = b.AvgSpeed.FormatValue("%.2f WPS")
	b.drawKeyVal(b.AvgSpeed.Key, val, keyStyle, valStyle)

	val = b.AvgAccuracy.FormatValue("%.2f %%")
	b.drawKeyVal(b.AvgAccuracy.Key, val, keyStyle, valStyle)
}
func (b *StatsBuffer) drawKeyVal(key, val string, keyStyle, valStyle tcell.Style) {
	b.Position.Y += 1
	keyX := b.Position.X + 2
	valX := b.Position.X + b.Size.Width - len(val) - 1
	screen.DrawText(b.screen, key, &keyX, &b.Position.Y, keyStyle)
	screen.DrawText(b.screen, val, &valX, &b.Position.Y, valStyle)
}
func (_ StatsBuffer) GetID() int {
	return STATS_BUFFER_ID
}
func (b *StatsBuffer) SetScreen(s tcell.Screen) {
	b.screen = s
}
func (b *StatsBuffer) HandleKeyEvent(_ *tcell.EventKey) {

}
func (b *StatsBuffer) SetWindow(w *screen.Window) {
	b.window = w
}

func (b *StatsBuffer) Update() {
	testsCount, err := b.Repository.CountAllWhere("")
	b.TestsCount = Stat{
		Key:   "Tests Completed:",
		Value: -1,
	}
	if err == nil {
		b.TestsCount.Value = testsCount
	}
	// all time time spent typing

	b.BestSpeed = Stat{
		Key:   "Best Speed:",
		Value: -1,
	}
	bestSpeed, err := b.Repository.MaxWhere("speed", "")
	if err == nil {
		bestSpeed := bestSpeed.(float64)
		b.BestSpeed.Value = float32(bestSpeed)
	}

	b.BestAccuracy = Stat{
		Key:   "Best Accuracy:",
		Value: -1,
	}
	bestAccuracy, err := b.Repository.MaxWhere("accuracy", "")
	if err == nil {
		ba := bestAccuracy.(float64)
		b.BestAccuracy.Value = float32(ba)
	}

	b.WordsCount = Stat{
		Key:   "Words Typed:",
		Value: -1,
	}
	wordsTyped, err := b.Repository.Sum("words_count", "")
	if err == nil {
		count := wordsTyped.(int64)
		b.WordsCount.Value = count
	}

	b.AvgSpeed = Stat{
		Key:   "Average Speed: ",
		Value: -1,
	}
	averageSpeed, err := b.Repository.Avg("speed", "")
	if err == nil {
		averageSpeed := averageSpeed.(float64)
		b.AvgSpeed.Value = averageSpeed
	}

	b.AvgAccuracy = Stat{
		Key:   "Average Accuracy: ",
		Value: -1,
	}
	avgAcc, err := b.Repository.Avg("accuracy", "")
	if err == nil {
		avgAcc := avgAcc.(float64)
		b.AvgAccuracy.Value = avgAcc
	}
}
