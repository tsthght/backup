package task

import (
	"sync"
	"time"
	"fmt"

	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/machine"
	"github.com/tsthght/backup/utils"
)

func Task(quit <-chan time.Time, wg *sync.WaitGroup, rate int, cluster *database.MGRInfo, user database.UserInfo) {
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
			//获取任务类型和任务状态，设置状态各个部分用 协程
			uuid := -1
			tp := ""
			var err error
			db := database.GetMGRConnection(cluster, user, true)
			if db == nil {
				fmt.Printf("db is nil")
				continue
			} else {
				uuid, err = database.GetTaskUUIDAsignedToMachine(db, ip)
				if err != nil {
					fmt.Printf("GetTaskUUIDAsignedToMachine failed: " + err.Error())
				}

				if uuid < 0 {
					fmt.Printf("no task todo now\n")
					continue
				}

				tp, err = database.GetTaskTypeByUUID(db, uuid)
				if err != nil {
					fmt.Printf("GetTaskTypeByUUID failed: " + err.Error())
				}
				db.Close()
			}

			switch tp {
			case "schema":
				fmt.Printf("do schema logic\n")
				go machine.StateMachineSchema()
			case "full":
				fmt.Printf("do full logic\n")
			case "all":
				fmt.Printf("do all logic\n")
			default:
				fmt.Printf("type is error\n")
			}
		}
	}
}