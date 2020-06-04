package machine

import (
	"database/sql"
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
			//更新状态
			initState = PrepareEnv
			database.SetMachineStageByIp(db, ip, "prepare_env")
			time.Sleep(2 * time.Second)
			//todo
		case PrepareEnv:
			//更新状态
			initState = PreCheck
			database.SetMachineStageByIp(db, ip, "pre_check")
			time.Sleep(2 * time.Second)
			//todo
		case PreCheck:
			//更新状态
			initState = Dumping
			database.SetMachineStageByIp(db, ip, "dumping")
			time.Sleep(2 * time.Second)
			//todo
		case Dumping:
			//更新状态
			initState = Loading
			database.SetMachineStageByIp(db, ip, "loading")
			time.Sleep(2 * time.Second)
			//todo
		case Loading:
			//更新状态
			initState = PosCheck
			database.SetMachineStageByIp(db, ip, "pos_check")
			time.Sleep(2 * time.Second)
			//todo
		case PosCheck:
			//更新状态
			initState = ResetEnv
			database.SetMachineStageByIp(db, ip, "reset_env")
			time.Sleep(2 * time.Second)
			//todo
		case ResetEnv:
			//更新状态
			initState = Done
			database.SetMachineStageByIp(db, ip, "done")
			time.Sleep(2 * time.Second)
			//todo
		case Done:
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
