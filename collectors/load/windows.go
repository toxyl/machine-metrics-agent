//go:build windows

package load

// stub because windows has no equivalent to 1m, 5m and 15m load average
func (la *Info) update() {}
