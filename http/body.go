package http

type baseInfo struct {
	Environment string           `json: 'env'`
	PhysicalClusterName string   `json: 'physical_cluster_name'`
	Username string              `json: 'username'`
}

type BinLogEnableInfo struct {
	baseInfo
	EnableBinlog bool            `json: 'enable_binlog'`
}

type PumpInfo struct {
	baseInfo
	Command string               `json: 'command'`
	Pumplist []string            `json: 'pumplist'`
}

type DrainerInfo struct {
	baseInfo
	Command string               `json: 'command'`
	Drainerlist []string            `json: 'drainerlist'`
}

type RollbackInfo struct {
	baseInfo
	Command string               `json: 'command'`
}