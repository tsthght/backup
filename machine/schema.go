package machine

import (
	"fmt"
	"time"

	"github.com/tsthght/backup/database"
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
)

func StateMachineSchema(cluster *database.MGRInfo, user database.UserInfo, initState int, ip string, uuid int) {
	loop:
		fmt.Printf("schema loop...\n")
		switch initState {
		case ToDo :
			fmt.Printf("state: todo\n")
			//更新状态
			initState = PrepareEnv
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "prepare_env")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "prepare_env")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			fmt.Printf("current state: %d\n", initState)
			//todo
			goto loop
		case PrepareEnv:
			fmt.Printf("state: prepare_env\n")
			//更新状态
			initState = PreCheck
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "pre_check")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "pre_check")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			fmt.Printf("current state: %d\n", initState)
			//todo
			goto loop
		case PreCheck:
			fmt.Printf("state: pre_check\n")
			//更新状态
			initState = Dumping
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "dumping")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "dumping")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			//todo
			goto loop
		case Dumping:
			fmt.Printf("state: dumping\n")
			//更新状态
			initState = Loading
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "loading")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "loading")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			//todo
			goto loop
		case Loading:
			fmt.Printf("state: loading\n")
			//更新状态
			initState = PosCheck
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "pos_check")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "pos_check")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			//todo
			goto loop
		case PosCheck:
			fmt.Printf("state: pos_check\n")
			//更新状态
			initState = ResetEnv
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "reset_env")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "reset_env")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			//todo
			goto loop
		case ResetEnv:
			fmt.Printf("state: reset_env\n")
			//更新状态
			initState = Done
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "done")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "done")
			}
			db.Close()
			time.Sleep(2 * time.Second)
			//todo
			goto loop
		case Done:
			fmt.Printf("state: done\n")
			time.Sleep(2 * time.Second)
			//更新状态
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				//应该限制次数的
				goto loop
			}
			err := database.SetMachineStageByIp(db, ip, "idle")
			if err != nil {
				fmt.Printf("call SetMachineStageByIp(%s, %s) failed\n", ip, "idle")
			}
			err = database.SetTaskStageByUUID(db, uuid, "success")
			if err != nil {
				fmt.Printf("call SetTaskStageByUUID(%s, %s) failed\n", ip, "success")
			}
			//设置任务状态
			db.Close()
			return
		}
}
