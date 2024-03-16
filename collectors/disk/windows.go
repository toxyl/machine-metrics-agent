//go:build windows

package disk

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func (d *Info) update() {
	var total, avail, used, bTotal, bAvail, bUsed uint64

	for drive := 'A'; drive <= 'Z'; drive++ {
		rootPath := fmt.Sprintf("%c:\\", drive)
		rootPathName, err := syscall.UTF16PtrFromString(rootPath)
		if err != nil {
			fmt.Printf("Error getting root path name for drive %s: %v\n", string(drive), err)
			continue
		}
		driveType := windows.GetDriveType(rootPathName)
		if driveType == windows.DRIVE_FIXED { // Only consider fixed drives
			var freeBytesAvailableToCaller, totalNumberOfBytes, totalNumberOfFreeBytes uint64

			err := windows.GetDiskFreeSpaceEx(
				rootPathName,
				&freeBytesAvailableToCaller,
				&totalNumberOfBytes,
				&totalNumberOfFreeBytes,
			)
			if err != nil {
				fmt.Printf("Error getting disk space information for drive %s: %v\n", string(drive), err)
				continue
			}

			total = uint64(totalNumberOfBytes)
			avail = uint64(totalNumberOfFreeBytes)
			used = total - avail

			bTotal += total
			bAvail += avail
			bUsed += used
		}
	}

	d.total = bTotal
	d.avail = bAvail
	d.used = bUsed
}
