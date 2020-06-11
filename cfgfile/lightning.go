package cfgfile

import (
	"html/template"
	"os"
)

const (
	LightningConfigFile = "lightning.toml"
	LightningLogFile = "tidb-lightning.log"
)

const lightningconfigfile = `
[lightning]
level = "info"
file = "{{ .Logfile }}"
# Prometheus
pprof-port = 8289
region-concurrency = {{ .ThreadNum }}

[checkpoint]
enable = false

[tikv-importer]
backend = "tidb"
on-duplicate = "replace"

[mydumper]
data-source-dir = "{{ .BKPath }}"

[tidb]
host = "{{ .Host }}"
port = {{ .Port }}
status-port = 10080
user = "{{ .Username }}"
password = "{{ .Password }}"
`

type LightingConfig struct {
	ThreadNum int
	BKPath string
	Port int
	Username string
	Password string
	Host string
	Logfile string
}

func GenLightningConfigFile(path, bkpath, username, password, host string, threadnum, port int) error {
	lc := LightingConfig{
		ThreadNum: threadnum,
		BKPath:    bkpath,
		Port:      port,
		Username:  username,
		Password:  password,
		Host:      host,
		Logfile:   LightningLogFile,
	}
	f, err := os.Create(path + "/" + LightningConfigFile)
	if err != nil {
		return err
	}
	t := template.Must(template.New("lightning").Parse(lightningconfigfile))
	return t.Execute(f, lc)
}
