package api

import (
	"container/list"
	"sync"
	"time"
)

type point struct {
	when   time.Time
	traced int
}

type rateMonitor struct {
	mu     sync.Mutex
	total  int
	points *list.List
}

func (m *rateMonitor) addPoint(when time.Time, traced int) {
	m.mu.Lock()
	m.total += traced
	cutoff := when.Add(-time.Minute)
	m.points.PushBack(point{when, m.total})
	for {
		f := m.points.Front()
		if f.Value.(point).when.After(cutoff) {
			break
		}
		m.points.Remove(f)
	}
	m.mu.Unlock()
}

func (m *rateMonitor) rateHz() float64 {
	m.mu.Lock()
	if m.points.Len() < 2 {
		m.mu.Unlock()
		return 0
	}
	epoch := time.Now()
	n := float64(m.points.Len())
	var sumXY, sumX, sumY, sumXX float64
	for e := m.points.Front(); e != nil; e = e.Next() {
		pt := e.Value.(point)
		x := pt.when.Sub(epoch).Seconds()
		y := float64(pt.traced)
		sumXY += x * y
		sumX += x
		sumY += y
		sumXX += x * x
	}
	m.mu.Unlock()
	return (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
}
