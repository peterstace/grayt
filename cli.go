package grayt

import (
	"fmt"
	"time"
)

type cli struct {
	firstDisplay bool

	start time.Time

	lastUpdate    time.Time
	lastCompleted uint64

	throughputSmoothed float64 // Completed per second.
}

func newCLI() *cli {
	now := time.Now()
	return &cli{true, now, now, 0, 0.0}
}

func (c *cli) update(completed, total uint64) {

	now := time.Now()

	// Calculate progress.
	progress := float64(completed) / float64(total) * 100

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
	remaining := total - completed
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
		displayDuration(now.Sub(c.start)),
		progress, displayFloat64(c.throughputSmoothed),
		displayDuration(etaDuration),
	)
	c.firstDisplay = false

	c.lastUpdate = now
	c.lastCompleted = completed
}

func (c cli) done() {
	fmt.Printf("\nDone.\n")
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
