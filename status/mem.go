package status

import (
	"errors"

	"github.com/shirou/gopsutil/mem"
	"github.com/tsthght/backup/database"
)

var meminfo database.MEMInfo

func getMemInfo() (info *database.MEMInfo, err error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	if mem == nil {
		return nil, errors.New("VirtualMemory return is nil")
	}
	if meminfo.TotalSize == 0 {
		meminfo.TotalSize = mem.Total
	}
	meminfo.Available = mem.Available
	meminfo.UsedPercent = mem.UsedPercent
	return &meminfo, nil
}
