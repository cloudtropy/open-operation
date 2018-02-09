package ticker

import (
  "time"
  "errors"
)


type Ticker struct {
  C chan time.Time
  d time.Duration
  timer *time.Timer
  running bool
}

func (t *Ticker)StartTimer() {
  t.running = true
  if t.timer != nil {
    t.C <- time.Now()
  }
  t.timer = time.AfterFunc(t.d, t.StartTimer)
}

func (t *Ticker)StopTimer() {
  if t.timer != nil {
    t.timer.Stop()
    t.timer = nil
  }
  t.running = false
  t.DoNow()
}

func (t *Ticker)Running() bool {
  return t.running
}

func (t *Ticker)DoNow() {
  t.C <- time.Now()
}

func NewTicker(d time.Duration) *Ticker {
  if d <= 0 {
    panic(errors.New("non-positive interval for NewTicker"))
  }

  c := make(chan time.Time, 1)
  t := &Ticker{
    C: c,
    d: d,
  }
  t.StartTimer()

  return t
}

