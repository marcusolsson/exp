package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"github.com/marcusolsson/exp/typiskt/round"
	"github.com/nsf/termbox-go"
)

var (
	screenWidth  int
	screenHeight int

	ticker    *time.Ticker
	currRound *round.Round
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "", log.LstdFlags)

	currRound = round.New(string(b), logger)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	screenWidth, screenHeight = termbox.Size()

	redraw(currRound)

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				switch currRound.State {
				case round.RoundStarted:
					stopRound()
				default:
					break loop
				}
			case termbox.KeySpace:
				switch currRound.State {
				case round.RoundStopped:
					if ev.Ch == 'r' {
						restartRound()
					}
				case round.RoundWaiting:
					startRound()
					fallthrough
				default:
					if err := currRound.Next(); err != nil {
						break
					}
				}
			case termbox.KeyBackspace2:
				currRound.Undo()
			default:
				switch currRound.State {
				case round.RoundStopped:
					if ev.Ch == 'r' {
						restartRound()
					}
				case round.RoundWaiting:
					startRound()
					fallthrough
				default:
					currRound.Advance(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		redraw(currRound)
	}
}

func startRound() {
	currRound.Start()
	ticker = time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			redraw(currRound)
		}
	}()
}

func stopRound() {
	ticker.Stop()
	currRound.Stop()
}

func restartRound() {
	stopRound()
	currRound.Reset()
	startRound()
}

func drawString(x, y, offx int, fg, bg termbox.Attribute, msg string) (int, int) {
	for _, c := range msg {
		termbox.SetCell(x+offx, y, c, fg, bg)
		x++
	}

	return x, y
}

func drawWord(x, y, offx int, tok round.Word) (int, int) {
	if tok.Curr {
		return drawString(x, y, offx, termbox.ColorDefault, termbox.ColorBlack, tok.Text)
	}
	if !tok.Done {
		return drawString(x, y, offx, termbox.ColorDefault, termbox.ColorDefault, tok.Text)
	}
	if !tok.OK {
		return drawString(x, y, offx, termbox.ColorRed, termbox.ColorDefault, tok.Text)
	}

	return drawString(x, y, offx, termbox.ColorGreen, termbox.ColorDefault, tok.Text)
}

func drawWords(x, y int, offx int, words []round.Word, lim int) int {
	for _, t := range words {
		if len(t.Text)+x > lim {
			y++
			x = 0
		}
		x, y = drawWord(x, y, offx, t)
		x++
	}

	return y
}

func redraw(r *round.Round) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	cx := (screenWidth / 2)

	y := drawWords(0, 2, cx-40, r.Words, 80)

	drawString(0, y+2, cx-40, termbox.ColorDefault, termbox.ColorDefault, "> "+r.Scratch)

	switch currRound.State {
	case round.RoundStarted:
		drawString(0, y+4, cx-2, termbox.ColorMagenta, termbox.ColorDefault, fmt.Sprintf("%.0f", math.Floor(time.Since(r.StartedAt).Seconds())))

	case round.RoundStopped:

		drawString(0, y+4, cx-40, termbox.ColorDefault, termbox.ColorDefault,
			fmt.Sprintf("You finished in %.0f seconds with an accuracy of %.2f%%",
				math.Floor(r.Took.Seconds()),
				r.Accuracy()))

		drawString(0, y+5, cx-40, termbox.ColorDefault, termbox.ColorDefault,
			fmt.Sprintf("You typed %d keystrokes (%d were wrong) and made %d corrections.",
				len(r.Typed()),
				len(r.Mistyped()),
				len(r.Corrections())))

		drawString(0, y+8, cx-10, termbox.ColorYellow, termbox.ColorDefault, "Press 'r' to restart")
		drawString(0, y+10, cx-10, termbox.ColorMagenta, termbox.ColorDefault, "Press 'Esc' to quit")

	case round.RoundWaiting:
		drawString(0, y+4, cx-14, termbox.ColorYellow, termbox.ColorDefault, "Start typing to begin round")
		drawString(0, y+6, cx-10, termbox.ColorMagenta, termbox.ColorDefault, "Press 'Esc' to quit")
	}

	termbox.Flush()
}
