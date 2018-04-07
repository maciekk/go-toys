// Toy problem from:
//   http://www.chilton.com/~jimw/ballclk.html
//
// Motivated by:
//   https://www.reddit.com/r/golang/comments/7awz6d/whats_a_good_project_to_do_to_learn_go/dpdpq2z/

// NOTE: This is initial naive version, without any optimizations.

package main

import "fmt"

// To distinguish balls, we will assign each of them an identifier.
type Ball int

// NOTE: in physical clock there is an extra, PERMANENT ball on the Hours rail,
// which results in a +1 count, thus yielding hours range of [1, 12] rather
// than [0, 11]. This is only for interpretation; code below ignores this.
type ClockState struct {
	// Top to bottom
	Mins     []Ball // 0-4; 5th one causes tip over
	FiveMins []Ball // 0-11; 12th one causes tip over
	Hours    []Ball // 0-11; 12th one causes tip over
	Queue    []Ball
}

// Append the balls from 'src', but in reverse order.
func AppendReverse(dest *[]Ball, src []Ball) {
	for i := len(src) - 1; i >= 0; i-- {
		*dest = append(*dest, src[i])
	}
}

// Create a fresh new clock state, in initial configuration.
func NewClock(num_balls int) ClockState {
	queue := make([]Ball, num_balls)
	for i := range queue {
		queue[i] = Ball(i)
	}
	return ClockState{Queue: queue}
}

// Advance clock by 1 minute.
func AdvanceClock(clock *ClockState) {
	var b Ball

	// Remove ball from queue.
	b, clock.Queue = clock.Queue[0], clock.Queue[1:]

	// Add it to 1-minute rail.
	clock.Mins = append(clock.Mins, b)

	// Finish if no spill-over.
	if len(clock.Mins) < 5 {
		return
	}

	// Spill over 1-minute rail.
	AppendReverse(&clock.Queue, clock.Mins[:4])
	clock.FiveMins = append(clock.FiveMins, b)
	clock.Mins = clock.Mins[:0]

	// Finish if no spill-over.
	if len(clock.FiveMins) < 12 {
		return
	}

	// Spill over 5-minute rail.
	AppendReverse(&clock.Queue, clock.FiveMins[:11])
	clock.Hours = append(clock.Hours, b)
	clock.FiveMins = clock.FiveMins[:0]

	// Finish if no spill-over.
	if len(clock.Hours) < 12 {
		return
	}

	// Spill over hours rail.
	AppendReverse(&clock.Queue, clock.Hours[:11])
	clock.Queue = append(clock.Queue, b)
	clock.Hours = clock.Hours[:0]
}

// Advance clock by exactly 24h.
func AdvanceClockOneDay(clock *ClockState) {
	const minutesPerDay = 24 * 60
	for i := 0; i < minutesPerDay; i++ {
		AdvanceClock(clock)
	}
}

// Checks if the clock is in its initial configuration, as would be returned
// by NewClock().
func isInitialConfig(clock ClockState) bool {
	// All rails other than Queue must be empty.
	if len(clock.Mins) > 0 || len(clock.FiveMins) > 0 || len(clock.Hours) > 0 {
		return false
	}
	// Queue balls must be all in order.
	for i, b := range clock.Queue {
		if b != Ball(i) {
			return false
		}
	}
	return true
}

// How many days before the clock state returns to its initial configuration?
func DaysUntilRepeat(numBalls int) int {
	c := NewClock(numBalls)
	AdvanceClockOneDay(&c) // To avoid premature loop termination on day 0
	day := 1
	for ; !isInitialConfig(c); day++ {
		AdvanceClockOneDay(&c)
	}
	return day
}

func main() {
	for i := 27; i <= 127; i++ {
		fmt.Printf("%d balls cycle after %d days.\n", i, DaysUntilRepeat(i))
	}
}
