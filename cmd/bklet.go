package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/tsthght/backup/args"
	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/secret"
)

func main() {
	//参数解析
	arg := args.Arguments{}
	args.InitArgs(&arg)
	fmt.Printf("cfg = %s\n", *arg.Cfg)
	//配置文件解析
	var conf config.BkConfig
	if _, err := toml.DecodeFile(*arg.Cfg, &conf); err != nil {
		fmt.Printf("error\n")
	}
	fmt.Printf("%v\n", conf)
	//参数转化


	//启动任务
	username := secret.GetValueByeKey(conf.Cmdb.Appkey, conf.Cmdb.Username)
	fmt.Printf("username: %s\n", username)
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
