package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tsthght/backup/args"
	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/register"
	"github.com/tsthght/backup/secret"
	"github.com/tsthght/backup/status"
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
	userinfo := database.UserInfo{
		Username: secret.GetValueByeKey(conf.Cmdb.Appkey, conf.Cmdb.Username),
		Password: secret.GetValueByeKey(conf.Cmdb.Appkey, conf.Cmdb.Password),
		Port:     strconv.Itoa(conf.Cmdb.Port),
		Database: conf.Cmdb.Database,
	}
	fmt.Printf("%v\n", userinfo)
	mgrinfo := database.MGRInfo{
		Hosts:      strings.Split(conf.Cmdb.Host, ","),
		WriteIndex: 0,
	}
	fmt.Printf("%v\n", mgrinfo)

	//启动任务
	wg := sync.WaitGroup{}
	quit := make(chan time.Time)

	//keepalive 5s
	wg.Add(1)
	go register.Register(quit, &wg, 5000, &mgrinfo, userinfo)

	//status 3s
	wg.Add(1)
	go status.Status(quit, &wg, 3000, &mgrinfo, userinfo, conf.Task.Path)

	wg.Wait()
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
}
