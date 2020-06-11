package call

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/tsthght/backup/cfgfile"
	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/execute"
)

func getLightningBackendArgs(cfg config.BkConfig, host string) ([]string, error) {
	var args []string
	//config
	args = append(args, "-config")
	args = append(args, cfg.Task.Path + "/" + cfgfile.LightningConfigFile)

	//pd
	args = append(args, "-pd-urls")
	args = append(args, host)
	return args, nil
}

func CallLightning(cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig, uuid int) error {
	//port
	p, err := strconv.Atoi(user.Port)
	if err != nil {
		return err
	}
	//host
	db := database.GetMGRConnection(cluster, user, true)
	if db == nil {
		return errors.New("db is nil")
	}
	bi, err := database.GetCluserBasicInfo(db, uuid, cfg, database.DownStream)
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	if len(bi.Hosts) == 0 {
		return errors.New("sql is nil")
	}
	idx := rand.Intn(len(bi.Hosts) - 1)
	//gen file
	err = cfgfile.GenLightningConfigFile(cfg.Task.Path, cfg.Task.Path + "/" + cfgfile.DataDir,  bi.User, bi.Password, bi.Hosts[idx], cfg.Task.DefaultLoaderThread, p)
	if err != nil {
		return err
	}
	//gen args
	if len(bi.ROOT) == 0 {
		return errors.New("root is nil")
	}
	idx = rand.Intn(len(bi.ROOT))
	args, err := getLightningBackendArgs(cfg, bi.ROOT[idx])
	if err != nil {
		return err
	}
	fmt.Printf("lightning args: %v\n", args)
	//call lightning
	output, err := execute.ExecuteCommand(cfg.Task.Path, "tidb-lightning", args...)
	if err != nil {
		fmt.Printf("call tidb-lightning failed. err : %s\n", output)
		return err
	}
	return nil
}

func CleanLightning (cfg config.BkConfig) error {
	os.Remove(cfg.Task.Path + "/" + cfgfile.LightningConfigFile)
	os.Remove(cfg.Task.Path + "/" + cfgfile.LightningLogFile)
	return nil
}