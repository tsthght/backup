package machine

import (
	"errors"

	"github.com/tsthght/backup/database"
)

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