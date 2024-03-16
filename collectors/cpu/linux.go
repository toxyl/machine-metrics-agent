//go:build linux

package cpu

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

type CPUTime struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	Irq     uint64
	Softirq uint64
	Steal   uint64
	Guest   uint64
}

func (ct CPUTime) Total() uint64 {
	return uint64(math.Max(0.0, float64(ct.User+ct.Nice+ct.System+ct.Idle+ct.Iowait+ct.Irq+ct.Softirq+ct.Steal+ct.Guest)))
}

func (c *Info) getCPUTime() CPUTime {
	var cpuTime CPUTime
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Println("Error opening /proc/stat:", err)
		return cpuTime
	}
	defer file.Close()

	var cpu string
	_, err = fmt.Fscanf(file, "%s %d %d %d %d %d %d %d %d %d",
		&cpu, &cpuTime.User, &cpuTime.Nice, &cpuTime.System, &cpuTime.Idle,
		&cpuTime.Iowait, &cpuTime.Irq, &cpuTime.Softirq, &cpuTime.Steal, &cpuTime.Guest)

	if err != nil {
		fmt.Println("Error reading /proc/stat:", err)
		return cpuTime
	}

	return cpuTime
}

func (c *Info) update() {
	initialTime := c.getCPUTime()
	time.Sleep(1 * time.Second)
	afterTime := c.getCPUTime()
	totalDiff := float64((afterTime.Total() - initialTime.Total()) * 100)
	c.UsedPct = math.Max(0.0, (totalDiff-float64((afterTime.Idle-initialTime.Idle)*100))/totalDiff)
	c.Cores = uint(runtime.NumCPU())
}
