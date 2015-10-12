package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	ui "github.com/gizak/termui"
)

var tps int64 = 60

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	initial()
}

func initial() {
	go handleClose(ui.EventCh())
	tickListener()
}

func handleClose(ech <-chan ui.Event) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	step := tickCh()

	for {
		select {
		case <-sig:
			return
		case e := <-ech:
			if e.Type == ui.EventKey && (e.Ch == 'q' || e.Ch == 'Q') {
				os.Exit(0)
			}
			if e.Type == ui.EventKey && e.Key == ui.KeyCtrlC {
				os.Exit(0)
			}
		case <-step:
			continue
		}
	}
}

var tickChs = make([]chan time.Time, 0)

func tickCh() <-chan time.Time {
	out := make(chan time.Time)
	tickChs = append(tickChs, out)
	return out
}

func tickListener() {
	tkr := time.NewTicker(time.Second / time.Duration(tps))
	defer tkr.Stop()

	go func() {
		for {
			t := <-tkr.C
			for _, tch := range tickChs {
				go func(ch chan time.Time) {
					ch <- t
				}(tch)
			}
		}
	}()
}
