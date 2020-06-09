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

type BladeConfig struct {
	BladeAk string                `toml:"bladeappkey" json:"bladeappkey"`
	BladeUser string              `toml:"bladeuser" json:"bladeuser"`
	BladePort int                 `toml:"bladeport" json:"bladeport"`
}

type TASK struct {
	Path string                  `toml:"path" json:"path"`
	DefaultGCTime string         `toml:"default-gc" json:"default-gc"`
	DefaultPump int              `toml:"default-pump" json:"default-pump"`
}

type BkConfig struct {
	Api API
	Cmdb CMDB
	Task TASK
	Blade BladeConfig
}