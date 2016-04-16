package grayt

import (
	"fmt"
	"time"
)

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
