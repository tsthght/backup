package status

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/tsthght/backup/database"
)



var cpuinfo database.CPUInfo

func getCPUInfo() (info *database.CPUInfo, err error) {
	cpuinfo.Percent = ""
	if cpuinfo.PhysicCoreNum == 0 {
		cpuinfo.PhysicCoreNum, err = cpu.Counts(false)
		if err != nil {
			return nil, err
		}
	}
	if cpuinfo.LogicCoreNum == 0 {
		cpuinfo.LogicCoreNum, err = cpu.Counts(true)
		if err != nil {
			return nil, err
		}
	}
	p, err := cpu.Percent(time.Duration(time.Second), false)
	if err != nil {
		return nil, err
	}
	for _, v := range p {
		cpuinfo.Percent += fmt.Sprintf("%f ", v)
	}
	return &cpuinfo, nil
}
