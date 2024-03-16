package disk

import (
	"math"
)

type Info struct {
	used  uint64
	total uint64
	avail uint64
}

func NewInfo() *Info {
	return &Info{
		used:  0,
		total: 0,
		avail: 0,
	}
}

func (d *Info) Collect() map[string]interface{} {
	d.update()
	return map[string]interface{}{
		"avail":         d.avail,
		"used":          d.used,
		"total":         d.total,
		"avail_percent": math.Max(0.0, float64(d.avail)/float64(d.total)),
		"used_percent":  math.Max(0.0, float64(d.used)/float64(d.total)),
	}
}
