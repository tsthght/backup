package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/tsthght/backup/config"
)

func main() {
	fmt.Printf("i am bklet\n")
	var conf config.BkConfig
	if _, err := toml.DecodeFile("../config/config.toml", &conf); err != nil {
		fmt.Printf("error\n")
	}
	fmt.Printf("## %v\n", conf)
	/*
	err := http.SetBinglogEnable(
		"http://xxxxxx:8000/api/v1/cluster/conf_cluster_binlog",
		"product",
		"inf_blade_multiidc",
		"ght",
		true,
		)
	err.Error()

	 */
	/*
	wg := sync.WaitGroup{}
	wg.Add(1)
	quit := make(chan time.Time)
	go register.Register(quit, &wg, 1000)
	wg.Wait()

	 */
}
