package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/tsthght/backup/args"
	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/secret"
)

func main() {
	//参数解析
	arg := args.ClientArgs{}
	args.InitClientArgs(&arg)
	//配置文件解析
	var conf config.ClientConfig
	if _, err := toml.DecodeFile(*arg.CfgFile, &conf); err != nil {
		fmt.Printf("error\n")
	}

	//参数转化
	userinfo := database.UserInfo{
		Username: secret.GetValueByeKey(conf.Cmdb.Appkey, conf.Cmdb.Username),
		Password: secret.GetValueByeKey(conf.Cmdb.Appkey, conf.Cmdb.Password),
		Port:     strconv.Itoa(conf.Cmdb.Port),
		Database: conf.Cmdb.Database,
	}
	mgrinfo := database.MGRInfo{
		Hosts:      strings.Split(conf.Cmdb.Host, ","),
		WriteIndex: 0,
	}

	db := database.GetMGRConnection(&mgrinfo, userinfo, true)
	if db == nil {
		fmt.Printf("call GetMGRConnection failed, err : %s\n", errors.New("db is nil"))
		return
	}

	switch *arg.Operator {
	case "create": {
		if len(*arg.Src) == 0 || len(*arg.Dst) == 0 {
			fmt.Printf("src and dst should not be nil\n")
			return
		}
		err := database.SetATask(db, *arg.Src, *arg.Dst, *arg.Type, *arg.Db)
		if err != nil {
			fmt.Printf("call SetATask failed, err : %s\n", errors.New("db is nil"))
			return
		}

		uuid, err := database.GetLatestTask(db, *arg.Src, *arg.Dst, *arg.Type, *arg.Db)
		if err != nil {
			fmt.Printf("get task uuid failed. err : %s\n", err.Error())
			return
		}
		fmt.Printf("crate task success, uuid = %d\n", uuid)
		return
	}
	case "show": {
		if strings.EqualFold(*arg.Role, "task") {
			str, err := database.GetTastInfo(db, *arg.UUID)
			if err != nil {
				fmt.Printf("get task info failed. err : %s\n", err.Error())
				return
			}
			fmt.Printf("Task info:\n\t%s\n", str)
			return
		} else {
			str, err := database.GetMachineInfo(db, *arg.UUID)
			if err != nil {
				fmt.Printf("get machin info failed. err : %s\n", err.Error())
				return
			}
			fmt.Printf("Machine info:\n\t")
			for _, v := range str {
				fmt.Printf("\t" + "%s", v)
			}
		}
	}
	}


}
