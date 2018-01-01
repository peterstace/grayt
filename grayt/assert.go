package grayt

import (
	"fmt"
	db "runtime/debug"
)

const enableAssertions = true

func assertUnit(v Vector) {
	if !enableAssertions {
		return
	}
	n2 := v.LengthSq()
	if n2 < 0.999 || n2 > 1.001 {
		db.PrintStack()
		panic(fmt.Sprintf("vector is not unit: %v", v))
	}
}
