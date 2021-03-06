package status

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/utils"
)

func Status(quit <-chan time.Time, wg *sync.WaitGroup, rate int, cluster *database.MGRInfo, user database.UserInfo, path string) {
	defer wg.Done()

	checkTick := time.NewTicker(time.Duration(rate) * time.Millisecond)
	defer checkTick.Stop()

	var ip string
	for {
		var err error
		err, ip = utils.GetLocalIP()
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if len(ip) > 0 {
			break
		}
	}

	for {
		select {
		case <- quit:
			fmt.Printf("cancel goroutine by channel")
			return
		case <- checkTick.C:
			cpu, err := getCPUInfo()
			if err != nil {
				continue
			}
			mem, err := getMemInfo()
			if err != nil {
				continue
			}
			dsk, err := getDiskInfo(path)
			if err != nil {
				continue
			}
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
			} else {
				_, err := database.StatusUpdateToCmdb(db, ip, *cpu, *mem, *dsk, path)
				if err != nil {
					fmt.Printf("status failed: %s", err.Error())
				}
				db.Close()
			}
		}
	}
}
