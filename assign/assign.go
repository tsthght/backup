package assign

import (
	"sync"
	"time"
	"fmt"

	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/utils"
)

func AssignTask(quit <-chan time.Time, wg *sync.WaitGroup, rate int, cluster *database.MGRInfo, user database.UserInfo) {
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
			//获取任务（state：todo ），更新（update）状态
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
			} else {
				_, err := database.AssignFromCmdb(db, ip)
				if err != nil {
					fmt.Printf("assign failed: %s", err.Error())
				}
				db.Close()
			}
		}
	}
}
