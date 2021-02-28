package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	runLoop()
}

func runLoop() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	render()
loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc,
					termbox.KeyCtrlC:
					break loop
				}
			}
		default:
			render()
			time.Sleep(250 * time.Millisecond)
		}
	}
	termbox.Close()
}

const coldef = termbox.ColorDefault

func render() {
	termbox.Clear(coldef, coldef)
	runPodList()
	termbox.Flush()
}
