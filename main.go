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
	if len(os.Args) < 2 {
		fmt.Printf("Usage:   %s [config]\n", os.Args[0])
		fmt.Printf("Example: %s config.yaml\n", os.Args[0])
		return
	}

	config, err := loadConfig(os.Args[1])
	if err != nil {
		panic("Error loading configuration: " + err.Error())
	}

	client := influx.NewClient(config.URL, config.Org, config.Bucket, config.Token, config.VerifyTLS)
	updateInterval := time.Duration(config.Interval) * time.Second
	hostname := "N/A"
	for {
		if h, err := os.Hostname(); err == nil {
			hostname = h
			if utils.IsRunningInLXC() {
				hostname += " (LXC)"
			}
		}

		t := time.Now()
		res := collectors.NewCollectorResults()
		wg := &sync.WaitGroup{}
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
