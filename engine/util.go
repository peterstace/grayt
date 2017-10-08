package engine

import "fmt"

const debug = true

func assertUnit(v vect3) {
	if !debug {
		return
	}
	n2 := v.norm2()
	if n2 <= 0.999 || n2 > 1.001 {
		panic(fmt.Sprintf("vector is not unit: norm2=%f vect3=%v", n2, v))
	}
}
