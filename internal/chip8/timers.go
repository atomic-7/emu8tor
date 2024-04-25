package chip8

import (
	"context"
	"time"
)

type ChipTimer struct {
	Count  uint8
	ticker *time.Ticker
}

func NewChipTimer() *ChipTimer {
	t := time.NewTicker(time.Second / 60)
	t.Stop()
	return &ChipTimer{Count: 0, ticker: t}
}

func (ct *ChipTimer) SetTimer(ctx context.Context, count uint8) {
	if ct.Count > 0 {
		ct.ticker.Reset(time.Second / 60)
	}
}

func (ct *ChipTimer) SetBeep(ctx context.Context, count uint8, beep func()) {
	go ct.decrease(ctx, beep)
}

func SilentBeeper() {
	println("Beep!")
}

func (ct *ChipTimer) decrease(ctx context.Context, beep func()) {
	select {
	case _ = <-ctx.Done():
		return
	case <-ct.ticker.C:
		if ct.Count > 0 {
			if beep != nil {
				beep()
			}
			ct.Count -= 1
		}
	}
}
