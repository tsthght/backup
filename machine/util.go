package machine

import (
	"errors"
	"fmt"

	"github.com/tsthght/backup/config"
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
	Failed

	Pump
	CheckPump
	Drainer
	CheckDrainer
)

var BKState map[string]int

func InitBKState() {
	BKState["todo"] = ToDo
	BKState["prepare_env"] = PrepareEnv
	BKState["pre_check"] = PreCheck
	BKState["dumping"] = Dumping
	BKState["loading"] = Loading
	BKState["pos_check"] = PosCheck
	BKState["reset_env"] = ResetEnv
	BKState["done"] = Done
	BKState["failed"] = Failed
	BKState["pump"] = Pump
	BKState["check_pump"] = CheckPump
	BKState["drainer"] = Drainer
	BKState["check_drainer"] = CheckDrainer
}


const (
	BKPATH = "bk"
)

func SetMachineStateByIp(cluster *database.MGRInfo, user database.UserInfo, ip, state string) error {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil.")
	}
	err := database.SetMachineStageByIp(db, ip, state)
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	return nil
}

func SetTaskState(cluster *database.MGRInfo, user database.UserInfo, uuid int, state, message string) error {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil")
	}
	err := database.SetTaskStateAndMessageByUUID(db, uuid, state, message)
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	return nil
}

func SetClusterGC(cluster *database.MGRInfo, user database.UserInfo, uuid int, cfg config.BkConfig, gc string) error {
	db := database.GetMGRConnection(cluster, user, false)
	if db == nil {
		return errors.New("db is nil")
	}

	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.UpStream)
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	db = database.GetTiDBConnection(bi)
	if db == nil {
		return errors.New("db is nil")
	}
	if len(gc) > 0 {
		err = database.SetGCTimeByUUID(db, gc)
	} else {
		err = database.SetGCTimeByUUID(db, cfg.Task.DefaultGCTime)
	}
	return err
}

func GetClusterGC(cluster *database.MGRInfo, user database.UserInfo, uuid int, cfg config.BkConfig) (error, string) {
	db := database.GetMGRConnection(cluster, user, false)
	if db == nil {
		return errors.New("mysql is nil"), ""
	}

	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.UpStream)
	if err != nil {
		db.Close()
		return err, ""
	}
	db.Close()
	fmt.Printf("bi: %v\n", bi)
	db = database.GetTiDBConnection(bi)
	if db == nil {
		return errors.New("tidb is nil"), ""
	}
	return database.GetGCTimeByUUID(db)
}