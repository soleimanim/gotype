package main

import (
	"log"
	"os"

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

	statusBuffer := buffers.NewStatusLineBuffer()
	actionsBuffer := buffers.NewActionsLineBuffer()
	typingTestBuffer := buffers.NewTypingTestBuffer(buffers.TestMode25Words)

	window := screen.NewWindow()
	err := window.Init(&statusBuffer)
	if err != nil {
		panic(err)
	}
	defer window.Close()

	window.AppendBuffer(actionsBuffer)
	window.AppendBuffer(&typingTestBuffer)

	window.Draw()

	for {
		if window.HandleEvents() {
			return
		}

		window.Draw()
	}
}

func handleCommandLineArgs() {

}
