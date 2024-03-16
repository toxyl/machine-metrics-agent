package collectors

type (
	CollectorFunc func() map[string]interface{}
	Collector     struct {
		Name string
		Host string
		fn   CollectorFunc
	}
	MetricGroup struct {
		Measurement string
		Host        string
		Fields      map[string]interface{}
	}
)

func (c *Collector) Collect() MetricGroup {
	return MetricGroup{
		Measurement: c.Name,
		Host:        c.Host,
		Fields:      c.fn(),
	}
}

func NewCollector(name, host string, fn CollectorFunc) *Collector {
	c := &Collector{
		Name: name,
		Host: host,
		fn:   fn,
	}
	return c
}
