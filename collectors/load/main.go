package load

import (
	"math"
	"runtime"
)

type Info struct {
	avg1  float64
	avg5  float64
	avg15 float64
}

func NewInfo() *Info {
	la := &Info{
		avg1:  0,
		avg5:  0,
		avg15: 0,
	}
	return la
}

func (l *Info) Collect() map[string]interface{} {
	l.update()
	cores := float64(runtime.NumCPU())
	return map[string]interface{}{
		"avg1m":          l.avg1,
		"avg5m":          l.avg5,
		"avg15m":         l.avg15,
		"avg1m_percent":  math.Max(0.0, float64(l.avg1)/cores),
		"avg5m_percent":  math.Max(0.0, float64(l.avg5)/cores),
		"avg15m_percent": math.Max(0.0, float64(l.avg15)/cores),
	}
}
