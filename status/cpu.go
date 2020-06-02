package status

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CPUInfo struct {
	LogicCoreNum int
	PhysicCoreNum int
	Percent string
}

var cpuinfo CPUInfo

func getCPUInfo() (info *CPUInfo, err error) {
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
	return &cpuinfo, nil
	for _, v := range p {
		cpuinfo.Percent += fmt.Sprintf("%f ", v)
	}
	return &cpuinfo, nil
}
