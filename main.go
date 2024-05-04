package main

import (
	"log"
	"os"

	"github.com/soleimanim/gotype/db"
	"github.com/soleimanim/gotype/logger"
	"github.com/soleimanim/gotype/screen"
	"github.com/soleimanim/gotype/screen/buffers"
)

func main() {
	args := os.Args

	for i, arg := range args {
		switch arg {
		case "--logFile":
			if i == len(args)-1 {
				log.Fatal("Error: The '--logFile' option cannot be used in isolation. Please provide the file path after it.")
			}
			path := args[i+1]
			file, err := logger.SetLoggerFile(path)
			if err != nil {
				log.Fatal("Error:", err)
			}
			defer file.Close()
		}
	}

	// Connect to database
	dbHandler := db.NewDBHandler()
	err := dbHandler.Open()
	if err != nil {
		// TODO: handle the error
		panic(err)
	}
	defer dbHandler.Close()

	typingTestRepository := db.NewTypingTestRepository(dbHandler.DB)
	// statusBuffer := buffers.NewStatusLineBuffer()
	// actionsBuffer := buffers.NewActionsLineBuffer(typingTestRepository)

	window := screen.NewWindow()
	err = window.Init()
	if err != nil {
		panic(err)
	}
	defer window.Close()

	// window.AppendBuffer(actionsBuffer)

	screenWidth, screenHeight := window.Screen.Size()
	recentTestBuffer := buffers.NewRecentTestsBuffer(screen.BufferPosition{
		X: 0,
		Y: 0,
	}, screen.BufferSize{
		Width:  screenWidth / 4,
		Height: screenHeight / 2,
	}, typingTestRepository)
	recentTestBuffer.Update()
	window.AppendBuffer(&recentTestBuffer)

	statsBuffer := buffers.NewStatsBuffer(screen.BufferPosition{
		X: 0,
		Y: screenHeight / 2,
	}, screen.BufferSize{
		Width:  screenWidth / 4,
		Height: screenHeight / 2,
	}, typingTestRepository)
	statsBuffer.Update()
	window.AppendBuffer(&statsBuffer)

	ttbPos, ttbSize := buffers.GetTypingTestBufferPositionAndSize(window.Screen)
	typingTestBuffer := buffers.NewTypingTestBuffer(ttbPos, ttbSize, buffers.TestMode25Words, typingTestRepository)
	window.AppendBuffer(&typingTestBuffer)
	window.Draw()

	for {
		if window.HandleEvents() {
			return
		}

		window.Draw()
	}
}
