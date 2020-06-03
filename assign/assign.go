package assign

import (
	"strings"
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
			//查看自己的状态
			db1 := database.GetMGRConnection(cluster, user, false)
			if db1 == nil {
				fmt.Printf("db1 is nil")
			} else {
				stage, err := database.GetStatusFromCmdb(db1, ip)
				if err != nil {
					fmt.Printf("get status failed: %s", err.Error())
				}
				fmt.Printf("## %s\n", stage)
				if !strings.EqualFold("idle", stage) {
					db1.Close()
					continue
				}
				db1.Close()
			}
			fmt.Printf("##########\n")

			//获取任务（state：todo ），更新（update）状态
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
			} else {
				uuid, err := database.AssignFromCmdb(db, ip)
				fmt.Printf("###### uuid=%d\n", uuid)
				if err != nil {
					fmt.Printf("assign failed: %s", err.Error())
				}
				db.Close()
			}
		}
	}
}
