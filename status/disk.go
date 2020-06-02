package status

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/tsthght/backup/database"
)

var diskinfo database.DiskInfo

func getDiskInfo(path string) (info *database.DiskInfo, err error) {
	dsk, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	diskinfo.TotalSize = dsk.Total
	diskinfo.Free = dsk.Free
	diskinfo.UsedPercent = dsk.UsedPercent
	return &diskinfo, nil
}
