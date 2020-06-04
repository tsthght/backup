package machine

import (
	"database/sql"
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



func StateMachineSchema(initState int, db *sql.DB, ip string, uuid int) {
	for {
		switch initState {
		case ToDo :
			fmt.Printf("state: todo\n")
			//更新状态
			initState = PrepareEnv
			database.SetMachineStageByIp(db, ip, "prepare_env")
			time.Sleep(2 * time.Second)
			//todo
		case PrepareEnv:
			fmt.Printf("state: prepare_env\n")
			//更新状态
			initState = PreCheck
			database.SetMachineStageByIp(db, ip, "pre_check")
			time.Sleep(2 * time.Second)
			//todo
		case PreCheck:
			fmt.Printf("state: pre_check\n")
			//更新状态
			initState = Dumping
			database.SetMachineStageByIp(db, ip, "dumping")
			time.Sleep(2 * time.Second)
			//todo
		case Dumping:
			fmt.Printf("state: dumping\n")
			//更新状态
			initState = Loading
			database.SetMachineStageByIp(db, ip, "loading")
			time.Sleep(2 * time.Second)
			//todo
		case Loading:
			fmt.Printf("state: loading\n")
			//更新状态
			initState = PosCheck
			database.SetMachineStageByIp(db, ip, "pos_check")
			time.Sleep(2 * time.Second)
			//todo
		case PosCheck:
			fmt.Printf("state: pos_check\n")
			//更新状态
			initState = ResetEnv
			database.SetMachineStageByIp(db, ip, "reset_env")
			time.Sleep(2 * time.Second)
			//todo
		case ResetEnv:
			fmt.Printf("state: reset_env\n")
			//更新状态
			initState = Done
			database.SetMachineStageByIp(db, ip, "done")
			time.Sleep(2 * time.Second)
			//todo
		case Done:
			fmt.Printf("state: done\n")
			time.Sleep(2 * time.Second)
			//更新状态
			database.SetMachineStageByIp(db, ip, "idle")
			database.SetTaskStageByUUID(db, uuid, "success")
			//设置任务状态
			db.Close()
			break
		}
	}
}
