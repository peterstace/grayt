package grayt

import (
	"fmt"
	"sync/atomic"
	"time"
)

func newProgress(total uint64) *progress {
	p := &progress{
		start:     time.Now(),
		completed: 0,
		total:     total,
		doneCh:    make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-p.doneCh:
				return
			case <-time.After(20 * time.Millisecond):

				pct := 100 * float64(atomic.LoadUint64(&p.completed)) / float64(p.total)
				samplesPerSecond := uint64(float64(p.completed) / time.Now().Sub(p.start).Seconds())
				remaining := time.Duration((p.total-p.completed)/samplesPerSecond) * time.Second

				fmt.Printf("\033[1K") // Erase from current post until start of line.
				fmt.Printf("\033[0E") // Move the cursor to the start of the current line.
				fmt.Printf("[ %5.2f%% ][ %.2e samples/sec ] Remaining: %s",
					pct, float64(samplesPerSecond), remaining)

			}
		}
	}()
	return p
}

type progress struct {
	start     time.Time
	completed uint64
	total     uint64
	doneCh    chan struct{}
}

func (b *progress) step() {
	atomic.AddUint64(&b.completed, 1)
}

func (b *progress) done() {
	close(b.doneCh)
	fmt.Println()
}
