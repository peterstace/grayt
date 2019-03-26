package api

import (
	"log"
	"sync"
	"time"
)

type render struct {
	scene   string
	pxWide  int
	pxHigh  int
	created time.Time

	cnd            *sync.Cond
	desiredWorkers int
}

func (r *render) work() {
	for {
		log.Println("...working...")
		time.Sleep(time.Second)
	}
}
