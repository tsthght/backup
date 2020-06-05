package database

const (
	UpStream = iota
	DownStream
)

type MGRInfo struct {
	Hosts []string
	WriteIndex int
}

type BladeInfo struct {
	Hosts []string
	User string
	Password string
	Port string
}

type UserInfo struct {
	Username string
	Password string
	Port string
	Database string
}

type CPUInfo struct {
	LogicCoreNum int
	PhysicCoreNum int
	Percent string
}

type MEMInfo struct {
	TotalSize uint64
	Available uint64
	UsedPercent float64
}

type DiskInfo struct {
	TotalSize uint64
	Free uint64
	UsedPercent float64
}