package machine

import (
	"errors"
	"fmt"
	"net"

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
	RollingSQL
	AddPump
	OpenBinlog
	Drainer
	CheckDrainer
	AddDrainer
	RollingMonitor
)

var BKState map[string]int

func InitBKState() {
	BKState = make(map[string]int)
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
	BKState["open_binlog"] = OpenBinlog
	BKState["drainer"] = Drainer
	BKState["check_drainer"] = CheckDrainer
	BKState["rolling_sql"] = RollingSQL
	BKState["add_pump"] = AddPump
	BKState["add_drainer"] = AddDrainer
	BKState["rolling_monitor"] = RollingMonitor
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

func SetTaskState(cluster *database.MGRInfo, user database.UserInfo, uuid int, state, stage, message string) error {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil")
	}
	err := database.SetTaskStateAndMessageByUUID(db, uuid, state, stage, message)
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
	err, str := database.GetGCTimeByUUID(db)
	db.Close()
	return err, str
}

func GetMachineNumByUUID(cluster *database.MGRInfo, user database.UserInfo, uuid int, stage string) (error, int) {
	db := database.GetMGRConnection(cluster, user, false)
	if db == nil {
		return errors.New("mysql is nil"), 0
	}
	err, num := database.GetMachineNum(db, uuid, stage)
	db.Close()
	return err, num
}

func GetLocalIP() (error, string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return err, ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return nil, ipnet.IP.String()
			}
		}
	}
	return errors.New("can not get local ip"), ""
}

func GetSrcClusterNameByUUID(cluster *database.MGRInfo, user database.UserInfo, uuid int) (error, string) {
	db := database.GetMGRConnection(cluster, user, false)
	if db == nil {
		return errors.New("mysql is nil"), ""
	}
	err, src := database.GetSrcClusterName(db, uuid)
	db.Close()
	return err, src
}

func GetMachinePumpIpByUUID(cluster *database.MGRInfo, user database.UserInfo, uuid int, stage string) (error, []string) {
	db := database.GetMGRConnection(cluster, user, false)
	if db == nil {
		return errors.New("mysql is nil"), nil
	}
	err, pumplist := database.GetMachinePumpIpByPump(db, uuid, stage)
	db.Close()
	return err, pumplist
}