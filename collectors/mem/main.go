package mem

import (
	"math"
)

type Info struct {
	total uint64
	avail uint64
}

func NewInfo() *Info {
	return &Info{}
}

func (m *Info) Collect() map[string]interface{} {
	m.update()
	fTotal := float64(m.total)
	return map[string]interface{}{
		"used":          m.total - m.avail,
		"used_percent":  math.Max(0.0, float64(m.total-m.avail)/fTotal),
		"avail":         m.avail,
		"avail_percent": math.Max(0.0, float64(m.avail)/fTotal),
		"total":         m.total,
	}
}
