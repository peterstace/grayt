package grayt

import (
	"fmt"
	"strings"
	"time"
)

type cli struct {
	pos   int
	total int

	lastDone      int
	lastUpdate    time.Time
	smoothedSpeed float64

	elapsed time.Duration
}

func newCLI(total int) *cli {
	return &cli{
		total: total,
	}
}

var posStrs = []string{`)`, `|`, `(`, `|`}

const updateInterval = 100 * time.Millisecond

func (c *cli) update(done int) {

	now := time.Now()
	timeDelta := now.Sub(c.lastUpdate)
	c.lastUpdate = now
	if timeDelta > 2*updateInterval {
		return
	} else {
		c.elapsed += timeDelta
	}

	doneDelta := done - c.lastDone
	c.lastDone = done

	if c.elapsed > 5*time.Second {
		speed := float64(doneDelta) / timeDelta.Seconds()
		const alpha = 0.001
		c.smoothedSpeed = alpha*speed + (1-alpha)*c.smoothedSpeed
	} else {
		elapsedInSeconds := float64(c.elapsed) / float64(time.Second)
		c.smoothedSpeed = float64(done) / elapsedInSeconds
	}

	eta := time.Duration(float64(c.total-done)/c.smoothedSpeed) * time.Second

	c.pos = (c.pos + 1) % len(posStrs)
	posStr := posStrs[c.pos]
	if done == c.total {
		posStr = strings.Repeat(" ", len(posStr))
	}

	pctDone := float64(done) / float64(c.total) * 100

	fmt.Print("\x1b[1G") // Move to column 1.
	fmt.Print("\x1b[2K") // Clear line.
	fmt.Printf(
		"%s [%6.2f%%] [Elapsed: %s] [ETA: %s] [%s samples/sec]",
		posStr,
		pctDone,
		displayDuration(c.elapsed),
		displayDuration(eta),
		displayFloat64(c.smoothedSpeed),
	)
}

func (c *cli) finished() {
	c.update(c.total)
	fmt.Printf("\n")
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
