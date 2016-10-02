package main

// (\) [XXX.XX%] [99.99Z samples/sec], ETA

import "fmt"

type cli struct {
	total int
}

func newCLI(total int) *cli {
	return &cli{total}
}

func (c *cli) update(done int) {
	fmt.Print("\x1b[1G") // Move to column 1.
	fmt.Print("\x1b[2K") // Clear line.
	fmt.Printf("%d/%d", done, c.total)
}

func (c *cli) finished() {
	fmt.Printf("\n")
}
