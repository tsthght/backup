package main

import (
	"fmt"

	"github.com/tsthght/backup/cfgfile"
)

func main() {
	err := cfgfile.GenLightningConfigFile("./lightning-backend.toml", "/data", "user", "pwd", "1.1.1.1", 32, 3306)
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
	}
}
