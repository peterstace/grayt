package grayt

import (
	"fmt"
	"sync/atomic"
	"time"
)

const cliUpdatePeriod = 100 * time.Millisecond

type cli struct {
	firstDisplay bool

	elapsed time.Duration

	lastUpdate    time.Time
	lastCompleted uint64

	throughputSmoothed float64 // Completed per second.

	total     uint64
	completed *uint64 // MUST only be read from atomically.

	done chan struct{}
}

func newCLI(total uint64, completed *uint64) *cli {
	now := time.Now()
	c := &cli{true, 0, now, 0, 0.0, total, completed, make(chan struct{})}
	go func() {
		for {
			var exit bool
			select {
			case <-c.done:
				exit = true
			case <-time.After(cliUpdatePeriod):
			}
			c.update()
			if exit {
				c.done <- struct{}{}
				return
			}
		}
	}()
	return c
}

func (c *cli) finish() {
	c.done <- struct{}{}
	<-c.done
}

func (c *cli) update() {

	now := time.Now()
	if now.Sub(c.lastUpdate) < 2*cliUpdatePeriod {
		c.elapsed += now.Sub(c.lastUpdate)
	}

	completed := atomic.LoadUint64(c.completed)

	// Calculate progress.
	progress := float64(completed) / float64(c.total) * 100

	// Calculate throughput.
	nowDelta := now.Sub(c.lastUpdate)
	completedDelta := completed - c.lastCompleted
	throughput := float64(completedDelta) / nowDelta.Seconds()
	const alpha = 0.001
	if c.throughputSmoothed == 0.0 {
		c.throughputSmoothed = throughput
	} else {
		c.throughputSmoothed = c.throughputSmoothed*(1.0-alpha) + throughput*alpha
	}

	// Calculate ETA.
	remaining := c.total - completed
	etaSec := float64(remaining) / c.throughputSmoothed
	etaDuration := time.Duration(etaSec*1e9) * time.Nanosecond

	// Display the output.
	if !c.firstDisplay {
		fmt.Print("\x1b[1G") // Move to column 1.
		for i := 0; i < 4; i++ {
			fmt.Print("\x1b[1A") // Move up.
			fmt.Print("\x1b[2K") // Clear line.
		}
	}
	fmt.Printf(
		"Elapsed:    %s\n"+
			"Progress:   %.2f%%\n"+
			"Throughput: %s samples/sec\n"+
			"ETA:        %s\n",
		displayDuration(c.elapsed),
		progress, displayFloat64(c.throughputSmoothed),
		displayDuration(etaDuration),
	)
	c.firstDisplay = false

	c.lastUpdate = now
	c.lastCompleted = completed
}

func displayFloat64(f float64) string {

	var thousands int

	for f >= 1000 {
		f /= 1000
		thousands++
	}

	suffix := [...]byte{' ', 'k', 'M', 'T', 'P', 'E'}[thousands]

	switch {
	case f < 10:
		return fmt.Sprintf("%.3f%c", f, suffix) // 9.999K
	case f < 100:
		return fmt.Sprintf("%.2f%c", f, suffix) // 99.99K
	case f < 1000:
		return fmt.Sprintf("%.1f%c", f, suffix) // 999.9K
	default:
		panic(f)
	}
}

func displayDuration(d time.Duration) string {
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute
	s := (d - h*time.Hour - m*time.Minute) / time.Second
	return fmt.Sprintf(
		"%d%d:%d%d:%d%d",
		h/10, h%10, m/10, m%10, s/10, s%10,
	)
}
