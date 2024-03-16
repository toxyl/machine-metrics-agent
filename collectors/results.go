package collectors

import (
	"sync"
	"time"

	"github.com/toxyl/machine-metrics-agent/influx"
)

type CollectorResults struct {
	mu      *sync.Mutex
	results map[string]MetricGroup
}

func (cr *CollectorResults) Set(key string, value MetricGroup) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.results[key] = value
}

func (cr *CollectorResults) WriteToInfluxDB(client influx.Client) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	t := time.Now()
	for _, v := range cr.results {
		client.WriteMultiple(v.Measurement, t, "host", v.Host, v.Fields)
	}
}

func NewCollectorResults() *CollectorResults {
	cr := &CollectorResults{
		mu:      &sync.Mutex{},
		results: map[string]MetricGroup{},
	}
	return cr
}
