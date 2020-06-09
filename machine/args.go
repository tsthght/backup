package machine

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
)

func PrepareDumpArgus(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, uuid int, tp int) (error, []string) {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil"), nil
	}

	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.UpStream)
	if err != nil {
		db.Close()
		return err, nil
	}

	var args []string = nil
	//host
	if len(bi.Hosts) == 0 {
		db.Close()
		return errors.New("hosts is nil"), nil
	} else {
		idx := rand.Intn(len(bi.Hosts))
		args = append(args, "-h")
		args = append(args, bi.Hosts[idx])
	}
	//user
	if len(bi.User) == 0 {
		db.Close()
		return errors.New("user is nil"), nil
	} else {
		args = append(args, "-u")
		args = append(args, bi.User)
	}

	//pwd
	if len(bi.Password) == 0 {
		db.Close()
		return errors.New("passowrd is nil"), nil
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
		db.Close()
		return err, nil
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

	if tp == 0 {
		//no data
		args = append(args, "-d")
	}
	db.Close()
	return nil, args
}

func PrepareLoadArgus(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, uuid int) (error, []string) {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil"), nil
	}

	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.DownStream)
	if err != nil {
		db.Close()
		return err, nil
	}

	var args []string = nil
	//host
	if len(bi.Hosts) == 0 {
		db.Close()
		return errors.New("hosts is nil"), nil
	} else {
		idx := rand.Intn(len(bi.Hosts))
		args = append(args, "-h")
		args = append(args, bi.Hosts[idx])
	}
	//user
	if len(bi.User) == 0 {
		db.Close()
		return errors.New("user is nil"), nil
	} else {
		args = append(args, "-u")
		args = append(args, bi.User)
	}

	//pwd
	if len(bi.Password) == 0 {
		db.Close()
		return errors.New("passowrd is nil"), nil
	} else {
		args = append(args, "-p")
		args = append(args, bi.Password)
	}

	//port
	args = append(args, "-P")
	args = append(args, bi.Port)

	//path
	args = append(args, "-d")
	args = append(args, BKPATH)

	db.Close()
	return nil, args
}

func PreparePumpArgus(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, uuid int) (error, []string) {
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil"), nil
	}

	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.DownStream)
	if err != nil {
		db.Close()
		return err, nil
	}
	db.Close()

	var args []string = nil
	//addr
	args = append(args, "-addr")
	args = append(args, "0.0.0.0:8250")
	//advertise-addr
	args = append(args, "-advertise-addr")
	err, ip := GetLocalIP()
	if err != nil {
		return err, nil
	}
	args = append(args, ip+ ":8250")
	//root
	args = append(args, "-pd-urls")
	var root []string
	for _, v := range bi.ROOT {
		root = append(root, "http://" + v + ":2379")
	}
	urls := strings.Join(root, ",")
	args = append(args, urls)
	//data-dir
	args = append(args, "-data-dir")
	args = append(args, "data.pump")
	//log-file
	args = append(args, "-log-file")
	args = append(args, "pump.log")
	//config
	args = append(args, "-config")
	args = append(args, "pump.toml")

	return nil, args
}