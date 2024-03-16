package uptime

type Info struct {
	up float64
}

func NewInfo() *Info {
	return &Info{
		up: 0,
	}
}

func (u *Info) Collect() map[string]interface{} {
	u.update()
	return map[string]interface{}{
		"seconds": u.up,
	}
}
