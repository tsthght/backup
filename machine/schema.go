package machine

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/execute"
)

const (
	ToDo = iota
	PrepareEnv
	PreCheck
	Dumping
	Loading
	PosCheck
	ResetEnv
	Done
	Failed
)

func StateMachineSchema(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, initState int, ip string, uuid int) {
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
			err := SetMachineStateByIp(cluster, user, ip, "pre_check")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			initState = PreCheck
		case PreCheck:
			//schema迁移，没必要修改gc时间
			err := SetMachineStateByIp(cluster, user, ip, "dumping")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			//更新状态
			initState = Dumping
		case Dumping:
			//调用dump
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				continue
			}

			bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.UpStream)
			if err != nil {
				fmt.Printf("call GetCluserBasicInfo failed.\n")
				db.Close()
				continue
			}

			var args []string = nil
			//host
			if len(bi.Hosts) == 0 {
				//应该报错
				continue
			} else {
				idx := rand.Intn(len(bi.Hosts))
				args = append(args, "-h")
				args = append(args, bi.Hosts[idx])
			}
			//user
			if len(bi.User) == 0 {
				//应该报错
				continue
			} else {
				args = append(args, "-u")
				args = append(args, bi.User)
			}

			//pwd
			if len(bi.Password) == 0 {
				//应该报错
				continue
			} else {
				args = append(args, "-p")
				args = append(args, bi.Password)
			}

			//port
			args = append(args, "-P")
			args = append(args, bi.Port)

			//db tb
			dbinfo, err := database.GetDBInfoByUUID(db, uuid)
			if err != nil {
				fmt.Printf("call GetDBInfoByUUID failed. error: %s\n", err.Error())
			}

			if dbinfo != "" {
				dbtb := strings.Split(dbinfo, ":")
				args = append(args, "-B")
				args = append(args, dbtb[0])
				if len(dbtb) == 2 && len(dbtb[1]) > 0 {
					args = append(args, "-T")
					args = append(args, dbtb[1])
				}
			}

			//path
			args = append(args, "-o")
			args = append(args, BKPATH)

			//no data
			args = append(args, "-d")

			db.Close()
			fmt.Printf("## bi= %v\n", bi)
			fmt.Printf("## before %v\n", time.Now())
			fmt.Printf("== rgs: as%v\n", args)
			output, err := execute.ExecuteCommand(cfg.Task.Path, "mydumper", args...)
			if err != nil || len(output) > 0{
				e := SetMachineStateByIp(cluster, user, ip, "failed")
				if e != nil {
					fmt.Printf("call SetMachineStateByIp failed. err : %s", e.Error())
					continue
				}
				initState = Failed
				e = SetTaskState(cluster, user, uuid, "failed", err.Error() + ";" + output)
				if e != nil {
					fmt.Printf("call SetTaskState faled. err : %s", e.Error())
					continue
				}
			}
			fmt.Printf("## after %v\n", time.Now())
			fmt.Printf("output: %s\n", string(output))

			e := SetMachineStateByIp(cluster, user, ip, "loading")
			if e != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", e.Error())
				continue
			}
			//更新状态
			initState = Loading
		case Loading:
			err := SetMachineStateByIp(cluster, user, ip, "pos_check")
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
			}
			initState = ResetEnv
			return
		case ResetEnv:
			//清理
			//os.RemoveAll(cfg.Task.Path + "/" + BKPATH)
			err := SetMachineStateByIp(cluster, user, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStateByIp failed. err : %s", err.Error())
			}
			// 直接返回
			return
		}
	}
}
