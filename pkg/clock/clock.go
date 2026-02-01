// Package clock provides time abstraction to facilitate tests
// and remove scattered time.Now() calls.
package clock

import "time"

// Clock abstracts time to allow fakes in tests.
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// Real returns the system clock.
var Real Clock = realClock{}

type realClock struct{}

func (realClock) Now() time.Time        { return time.Now() }
func (realClock) Sleep(d time.Duration) { time.Sleep(d) }

// Fake is a controllable clock for tests.
type Fake struct {
	NowVal time.Time
}

// Now returns the configured time.
func (f *Fake) Now() time.Time { return f.NowVal }

// Sleep advances NowVal by d (does not block).
func (f *Fake) Sleep(d time.Duration) { f.NowVal = f.NowVal.Add(d) }

// Advance advances the time by d.
func (f *Fake) Advance(d time.Duration) { f.NowVal = f.NowVal.Add(d) }
