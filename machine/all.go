package machine

import (
	"fmt"
	"os"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/execute"
)

func StateMachineAll(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, initState int, ip string, uuid int, tp int) {
	message := ""
	gctime := ""
	for {
		fmt.Printf("schema loop...\n")
		switch initState {
		case Pump:
			//修改任务状态为 (doing, pump)
			err := SetTaskState(cluster, user, uuid, "doing", "pump", "")
			if err != nil {
				fmt.Printf("call SetTaskState faled. err : %s", err.Error())
				continue
			}

			//判断是否需要启动pump
			err, num := GetMachineNumByUUID(cluster, user, uuid, "pump")
			if err != nil {
				fmt.Printf("call GetMachineNumByUUID failed. err : %s", err.Error())
				continue
			}
			if num >= cfg.Task.DefaultPump {
				//更新状态为 修改配置文件
				err := SetTaskState(cluster, user, uuid, "todo", "rolling_sql", "")
				if err != nil {
					fmt.Printf("call SetTaskState failed. err : %s", err.Error())
					continue
				}
				//更新机器状态，更新任务状态
				err = SetMachineStateByIp(cluster, user, ip, "idle")
				if err != nil {
					fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
					continue
				}
				return
			}

			//设置 machine状态为pump
			err = SetMachineStateByIp(cluster, user, ip, "pump")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}
			//启动pump
			err, args := PreparePumpArgus(cluster, user, cfg, uuid)
			if err != nil {
				fmt.Printf("call PreparePumpArgus failed. err : %s", err.Error())
				continue
			}
			//阻塞
			_, err = execute.ExecuteCommand(cfg.Task.Path, "pump", args...)
			if err != nil {
				//忽略错误
				fmt.Printf("call ExecuteCommand failed. error : %s", err.Error())
			}
			os.Remove(cfg.Task.Path + "/" + "pump.log")
			os.Remove(cfg.Task.Path + "/" + "data.pump")

			//修改自己状态
			err = SetMachineStateByIp(cluster, user, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}
		case OpenBinlog:
			//修改自己状态为openbinlog
			//修改任务状态(doing,open_binlog)
			//调用接口
			//周期性检查是否打开

			//确认打开后，更新状态
			initState = PreCheck
		default:
			fmt.Printf("state is error\n")
			return
		}
	}
}
