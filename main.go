package main

import (
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	onExit := func() {
	}

	systray.Run(onReady, onExit)
}

type Timer struct {
	start      time.Time
	ticker     *time.Ticker
	pausedDiff time.Duration
	paused     bool
}

func (t *Timer) runTimer() {
	t.ticker = time.NewTicker(time.Second)

	defer func() {
		t.ticker.Stop()
	}()

	if !t.paused {
		t.start = time.Now()
	} else {
		t.start = time.Now().Add(-t.pausedDiff)
	}
	t.paused = false
	for {
		select {
		case <-t.ticker.C:
			t.setTitle(beautifyTime(time.Now().Sub(t.start)))
		}
	}
}

func (t *Timer) resetTimer() {
	t.start = time.Now()
	t.setTitle("Trimer")
}

func (t *Timer) pauseTimer() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.paused = true
		t.pausedDiff = time.Now().Sub(t.start)
	}
}

func (t *Timer) stopTimer() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.paused = false
		t.pausedDiff = 0
		t.setTitle("Trimer")
	}
}

func (t *Timer) setTitle(s string) {
	n := len(s)
	const max = 12

	if n < max {
		var buf [max]byte
		isEven := (n % 2) == 0
		diff := max - n
		var padding string
		if isEven {
			padding = strings.Repeat(" ", diff/2)
		} else {
			padding = strings.Repeat(" ", (diff-1)/2)
		}
		m := len(padding)
		copy(buf[:m], padding)
		copy(buf[m:n+m], s)
		copy(buf[m+n:], padding)

		systray.SetTitle(string(buf[:]))
	} else {
		systray.SetTitle(s)
	}
}

func beautifyTime(d time.Duration) string {
	return d.Truncate(time.Second).String()
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Trimer")
	systray.SetTooltip("Trimer")
	startTimer := systray.AddMenuItem("Start", "Start")
	resetTimer := systray.AddMenuItem("Reset", "Reset the timer")
	pauseTimer := systray.AddMenuItem("Pause", "Pause the timer")
	stopTimer := systray.AddMenuItem("Stop", "Stop the timer")

	systray.AddSeparator()

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	timer := Timer{}

	// We can manipulate the systray in other goroutines
	go func() {
		for {
			select {
			case <-startTimer.ClickedCh:
				go timer.runTimer()
			case <-resetTimer.ClickedCh:
				timer.resetTimer()
			case <-pauseTimer.ClickedCh:
				timer.pauseTimer()
			case <-stopTimer.ClickedCh:
				timer.stopTimer()
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}
