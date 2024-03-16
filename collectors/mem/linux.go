//go:build linux

package mem

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (m *Info) update() *Info {
	const meminfoPath = "/proc/meminfo"

	data, err := os.ReadFile(meminfoPath)
	if err != nil {
		fmt.Println("Error reading meminfo:", err)
		return m
	}

	lines := strings.Split(string(data), "\n")
	meminfo := make(map[string]uint64)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		meminfo[key] = val
	}

	m.total = meminfo["MemTotal"] * 1024     // values are in kB
	m.avail = meminfo["MemAvailable"] * 1024 // values are in kB

	return m
}
