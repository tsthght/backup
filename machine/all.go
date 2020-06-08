package machine

import (
	"fmt"
	"os"
	"time"

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
		case ToDo :
			err := SetMachineStateByIp(cluster, user, ip, "prepare_env")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			initState = PrepareEnv
		case PrepareEnv:
			if tp != 0 {
				err, gc := GetClusterGC(cluster, user, uuid, cfg)
				if err != nil {
					fmt.Printf("call GetClusterGC failed. err : %s", err.Error())
					continue
				}
				gctime = gc
				//修改GC时间
				err = SetClusterGC(cluster, user, uuid, cfg, "168h")
				if err != nil {
					fmt.Printf("call SetClusterGC failed. err : %s", err.Error())
					continue
				}
				time.Sleep(10 * time.Second)
			}

			err := SetMachineStateByIp(cluster, user, ip, "pre_check")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			initState = PreCheck
		case PreCheck:
			err := SetMachineStateByIp(cluster, user, ip, "dumping")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			//更新状态
			initState = Dumping
		case Dumping:
			err, args := PrepareDumpArgus(cluster, user, cfg, uuid, tp)
			if err != nil {
				fmt.Printf("call PrepareDumpArgus failed. err : %s", err.Error())
				continue
			}

			output, err := execute.ExecuteCommand(cfg.Task.Path, "mydumper", args...)
			if err != nil || len(output) > 0{
				e := SetMachineStateByIp(cluster, user, ip, "failed")
				if e != nil {
					fmt.Printf("call SetMachineStateByIp failed. err : %s", e.Error())
					continue
				}
				message = ""
				if err != nil {
					message += err.Error()
				}
				if len(output) > 0 {
					message += output
				}
				initState = Failed
				continue
			}

			e := SetMachineStateByIp(cluster, user, ip, "loading")
			if e != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", e.Error())
				continue
			}
			//更新状态
			initState = Loading
		case Loading:
			err, args := PrepareLoadArgus(cluster, user, cfg, uuid)
			if err != nil {
				fmt.Printf("call PrepareLoadArgus failed. err : %s", err.Error())
				continue
			}
			output, err := execute.ExecuteCommand(cfg.Task.Path, "loader", args...)
			if err != nil {
				e := SetMachineStateByIp(cluster, user, ip, "failed")
				if e != nil {
					fmt.Printf("call SetMachineStateByIp failed. err : %s", e.Error())
					continue
				}
				message = ""
				if err != nil {
					message += err.Error()
				}
				if len(output) > 0 {
					message += output
				}
				initState = Failed
				continue
			}

			err = SetMachineStateByIp(cluster, user, ip, "pos_check")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			//更新状态
			initState = PosCheck
		case PosCheck:
			err := SetMachineStateByIp(cluster, user, ip, "done")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			//更新状态
			initState = Done
		case Done:
			err := SetMachineStateByIp(cluster, user, ip, "reset_env")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			err = SetTaskState(cluster, user, uuid, "success", "")
			if err != nil {
				fmt.Printf("call SetTaskState failed. err : %s", err.Error())
			}
			initState = ResetEnv
			continue
		case Failed:
			err := SetMachineStateByIp(cluster, user, ip, "reset_env")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}

			err = SetTaskState(cluster, user, uuid, "failed", message)
			if err != nil {
				fmt.Printf("call SetTaskState faled. err : %s", err.Error())
				continue
			}

			initState = ResetEnv
			continue
		case ResetEnv:
			//修改GC时间
			err := SetClusterGC(cluster, user, uuid, cfg, gctime)
			if err != nil {
				fmt.Printf("call SetClusterGC failed. err : %s", err.Error())
				continue
			}

			//清理
			os.RemoveAll(cfg.Task.Path + "/" + BKPATH)
			err = SetMachineStateByIp(cluster, user, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
				continue
			}
			// 直接返回
			return
		}
	}
}
