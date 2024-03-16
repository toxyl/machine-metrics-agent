package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/toxyl/machine-metrics-agent/collectors"
	"github.com/toxyl/machine-metrics-agent/collectors/cpu"
	"github.com/toxyl/machine-metrics-agent/collectors/disk"
	"github.com/toxyl/machine-metrics-agent/collectors/load"
	"github.com/toxyl/machine-metrics-agent/collectors/mem"
	"github.com/toxyl/machine-metrics-agent/collectors/net"
	"github.com/toxyl/machine-metrics-agent/collectors/uptime"
	"github.com/toxyl/machine-metrics-agent/influx"
	"github.com/toxyl/machine-metrics-agent/utils"
)

var (
	cpuInfo      = cpu.NewInfo()
	memInfo      = mem.NewInfo()
	loadInfo     = load.NewInfo()
	netInfo      = net.NewInfo()
	diskInfo     = disk.NewInfo()
	uptimeInfo   = uptime.NewInfo()
	collectorFns = map[string]collectors.CollectorFunc{
		"cpu":    cpuInfo.Collect,
		"mem":    memInfo.Collect,
		"load":   loadInfo.Collect,
		"disk":   diskInfo.Collect,
		"net":    netInfo.Collect,
		"uptime": uptimeInfo.Collect,
	}
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error retrieving hostname:", err)
		return
	}

	if utils.IsRunningInLXC() {
		hostname += " (LXC)"
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	client := influx.NewClient(config.URL, config.Org, config.Bucket, config.Token)

	updateInterval := time.Duration(config.Interval) * time.Second
	for {
		t := time.Now()
		wg := &sync.WaitGroup{}
		res := collectors.NewCollectorResults()
		for k, v := range collectorFns {
			wg.Add(1)
			go func() {
				defer wg.Done()
				res.Set(k, collectors.NewCollector(k, hostname, v).Collect())
			}()
		}
		wg.Wait()
		res.WriteToInfluxDB(*client)
		sleep := updateInterval - time.Since(t)
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}
