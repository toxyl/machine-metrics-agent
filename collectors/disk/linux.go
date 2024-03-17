//go:build linux

package disk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ignoreFS = []string{
		"autofs",
		"binfmt_misc",
		"bpf",
		"cgroup2",
		"configfs",
		"debugfs",
		"devpts",
		"devtmpfs",
		"efivarfs",
		"fuse.gvfsd-fuse",
		"fuse.portal",
		"fusectl",
		"hugetlbfs",
		"lxcfs",
		"mqueue",
		"none",
		"nsfs",
		"overlay",
		"proc",
		"pstore",
		"ramfs",
		"securityfs",
		"squashfs",
		"sunrpc",
		"sysfs",
		"tmpfs",
		"tracefs",
		"udev",
		"zfs",
	}
)

type MountpointData struct {
	ZFSPools           []ZFSPool
	MountedFilesystems []MountedFilesystem
	Total              TotalDiskData
}

type ZFSPool struct {
	Name  string
	Total int64
	Avail int64
	Used  int64
	Usage float64
}

type MountedFilesystem struct {
	Mountpoint string
	Total      int64
	Avail      int64
	Used       int64
	Usage      float64
}

type TotalDiskData struct {
	Total int64
	Avail int64
	Used  int64
}

func calculateZFSDiskUsage(dataset string) (int64, error) {
	output, err := exec.Command("zfs", "list", "-Hp", "-o", "used", dataset).Output()
	if err != nil {
		return 0, err
	}
	usedSpace, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}
	return usedSpace, nil
}

func calculateZFSTotalSize(dataset string) (int64, error) {
	output, err := exec.Command("zfs", "list", "-Hp", "-o", "used,avail", dataset).Output()
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(output))
	usedSpace, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	availSpace, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return usedSpace + availSpace, nil
}

func calculateDiskUsage(mountpoint string) (int64, int64, error) {
	output, err := exec.Command("df", "--output=used,avail", "-B", "1", mountpoint).Output()
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	fields := strings.Fields(lines[len(lines)-1])
	usedBlocks, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	availBlocks, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return usedBlocks, availBlocks, nil
}

func collectDiskData() MountpointData {
	var mountpointData MountpointData

	zpoolOutput, err := exec.Command("zpool", "list", "-H", "-o", "name").Output()
	if err == nil {
		zpoolLines := strings.Split(strings.TrimSpace(string(zpoolOutput)), "\n")
		for _, pool := range zpoolLines {
			if pool == "" {
				continue
			}
			var zfsPool ZFSPool
			zfsPool.Name = pool
			zfsDiskTotal, err := calculateZFSTotalSize(pool)
			if err != nil {
				fmt.Printf("Error calculating ZFS total size for pool %s: %v\n", pool, err)
				continue
			}
			zfsPool.Total = zfsDiskTotal
			zfsDiskUsed, err := calculateZFSDiskUsage(pool)
			if err != nil {
				fmt.Printf("Error calculating ZFS disk usage for pool %s: %v\n", pool, err)
				continue
			}
			zfsPool.Used = zfsDiskUsed
			zfsPool.Avail = zfsDiskTotal - zfsDiskUsed
			zfsPool.Usage = float64(zfsDiskUsed) * 100 / float64(zfsDiskTotal)
			mountpointData.ZFSPools = append(mountpointData.ZFSPools, zfsPool)
		}
	}

	file, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Println("Error opening /proc/mounts:", err)
		return mountpointData
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		mountpoint := fields[1]
		fsType := fields[2]

		skip := false
		for _, fs := range ignoreFS {
			if fsType == fs {
				skip = true
				break
			}
		}

		if !skip {
			var mountedFs MountedFilesystem
			mountedFs.Mountpoint = mountpoint
			diskUsed, diskAvail, err := calculateDiskUsage(mountpoint)
			if err != nil {
				fmt.Printf("Error calculating disk usage for mountpoint %s: %v\n", mountpoint, err)
				continue
			}
			if total := diskUsed + diskAvail; total != 0 {
				mountedFs.Total = total
				mountedFs.Used = diskUsed
				mountedFs.Avail = diskAvail
				mountedFs.Usage = float64(diskUsed) * 100 / float64(total)
				mountpointData.MountedFilesystems = append(mountpointData.MountedFilesystems, mountedFs)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading /proc/mounts:", err)
		return mountpointData
	}

	for _, pool := range mountpointData.ZFSPools {
		mountpointData.Total.Total += pool.Total
		mountpointData.Total.Avail += pool.Avail
		mountpointData.Total.Used += pool.Used
	}
	for _, fs := range mountpointData.MountedFilesystems {
		mountpointData.Total.Total += fs.Total
		mountpointData.Total.Avail += fs.Avail
		mountpointData.Total.Used += fs.Used
	}

	return mountpointData
}

func (d *Info) update() {
	stats := collectDiskData()
	d.total = uint64(stats.Total.Total)
	d.used = uint64(stats.Total.Used)
	d.avail = uint64(stats.Total.Avail)
}
