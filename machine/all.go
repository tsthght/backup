package machine

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/execute"
	"github.com/tsthght/backup/http"
)

func StateMachineAll(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, initState int, ip string, uuid int) {

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
			fmt.Printf("current num: %d, default num: %d\n", num, cfg.Task.DefaultPump)
			if num >= cfg.Task.DefaultPump {
				//更新状态为 修改配置文件
				err := SetTaskState(cluster, user, uuid, "todo", "add_pump", "")
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
			//设置 task的状态
			err = SetTaskState(cluster, user, uuid, "todo", "pump", "")
			if err != nil {
				fmt.Printf("call SetTaskState failed. err : %s", err.Error())
				continue
			}
			//启动pump
			err, args := PreparePumpArgus(cluster, user, cfg, uuid)
			if err != nil {
				fmt.Printf("call PreparePumpArgus failed. err : %s", err.Error())
				continue
			}
			//阻塞
			output, err := execute.ExecuteCommand(cfg.Task.Path, "pump", args...)
			if err != nil {
				//忽略错误
				fmt.Printf("call ExecuteCommand failed. error : %s, message : %s", err.Error(), output)
			}
			os.Remove(cfg.Task.Path + "/" + "pump.log")
			os.Remove(cfg.Task.Path + "/" + "data.pump")

			//修改自己状态
			err = SetMachineStateByIp(cluster, user, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}
		case AddPump:
			//修改任务状态为 (doing, pump)
			err := SetTaskState(cluster, user, uuid, "doing", "add_pump", "")
			if err != nil {
				fmt.Printf("call SetTaskState faled. err : %s", err.Error())
				continue
			}
			//获取src
			err, src := GetSrcClusterNameByUUID(cluster, user, uuid)
			if err != nil {
				fmt.Printf("call GetSrcClusterNameByUUID failed. err : %s\n", err.Error())
				continue
			}
			//获取pump地址
			err, pl := GetMachinePumpIpByUUID(cluster, user, uuid, "pump")
			if err != nil {
				fmt.Printf("call GetMachinePumpIpByUUID failed. err : %s\n", err.Error())
				continue
			}
			var pumplist []string
			for i, v := range pl {
				pumplist = append(pumplist, "pump" + strconv.Itoa(i) + " ansible_host=" + v + " deploy_dir="+ cfg.Task.Path)
			}
			err = http.SetPumpStatus(cfg.Api.ConfigPump, "product", src, "person", "append", pumplist)
			if err != nil {
				fmt.Printf("call SetPumpStatus faild. err : %s\n", err.Error())
				continue
			}
			time.Sleep(2 * time.Second)
			//设置 machine状态为pump
			err = SetMachineStateByIp(cluster, user, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}
			//设置 task的状态
			err = SetTaskState(cluster, user, uuid, "todo", "open_binlog", "")
			if err != nil {
				fmt.Printf("call SetTaskState failed. err : %s", err.Error())
				continue
			}
		case OpenBinlog:

		case Dumping:

		case Drainer:

		case AddDrainer:

		case RollingMonitor:

		default:
			fmt.Printf("state is error\n")
			return
		}
	}
}
