package cpu

type Info struct {
	UsedPct float64
	Cores   uint
}

func NewInfo() *Info {
	return &Info{
		UsedPct: 0,
		Cores:   0,
	}
}

func (c *Info) Collect() map[string]interface{} {
	c.update()
	return map[string]interface{}{
		"used_percent": c.UsedPct,
		"cores":        c.Cores,
	}
}
