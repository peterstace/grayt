package grayt

import (
	"fmt"
	db "runtime/debug"
)

func assertUnit(v Vector) {
	if !*debug {
		return
	}
	n2 := v.LengthSq()
	if n2 < 0.999 || n2 > 1.001 {
		db.PrintStack()
		panic(fmt.Sprintf("vector is not unit: %v", v))
	}
}
