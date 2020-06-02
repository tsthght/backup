package config

type API struct {
	ConfigBinlog string         `toml:"config_binlog" json:"config_binlog"`
	ConfigPump string           `toml:"config_pump" json:"config_pump"`
	ConfigDrainer string        `toml:"config_drainer" json:"config_drainer"`
	ConfigReset string          `toml:"config_restart" json:"config_restart"`
}

type CMDB struct {
	Appkey string                `toml:"appkey" json:"appkey"`
	Username string              `toml:"username" json:"username"`
	Password string              `toml:"password" json:"password"`
	Port int                     `toml:"cmdb_port" json:"cmdb_port"`
	Host string                  `toml:"cmdb_host" json:"cmdb_host"`
	Database string              `toml:"cmdb_db" json:"cmdb_db"`
}

type TASK struct {
	Path string
}

type BkConfig struct {
	Api API
	Cmdb CMDB
	Task TASK
}