package net

type Info struct {
	in  uint64
	out uint64
}

func NewInfo() *Info {
	return &Info{
		in:  0,
		out: 0,
	}
}

func (n *Info) Collect() map[string]interface{} {
	n.update()
	return map[string]interface{}{
		"in":  n.in,
		"out": n.out,
	}
}
