package main

import (
	"io/fs"
	"log"
	"os"

	"github.com/soleimanim/gotype/screen"
	"github.com/soleimanim/gotype/screen/buffers"
)

func main() {
	file, err := os.OpenFile("./log", os.O_APPEND|os.O_WRONLY, fs.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	statusBuffer := buffers.NewStatusLineBuffer()
	actionsBuffer := buffers.NewActionsLineBuffer()
	typingTestBuffer := buffers.NewTypingTestBuffer(buffers.TestMode25Words)

	window := screen.NewWindow()
	err = window.Init(&statusBuffer)
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
