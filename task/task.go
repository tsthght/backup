package task

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/database"
	"github.com/tsthght/backup/machine"
	"github.com/tsthght/backup/utils"
)

func Task(quit <-chan time.Time, wg *sync.WaitGroup, rate int, cluster *database.MGRInfo, user database.UserInfo, cfg config.BkConfig) {
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

	go func () {
		for {
			time.Sleep(time.Duration(rate + rand.Intn(500)) * time.Millisecond)
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
					db.Close()
					continue
				}

				if uuid < 0 {
					fmt.Printf("no task todo now\n")
					db.Close()
					continue
				}

				stage, err := database.GetMachineStageByIp(db, ip)
				if err != nil {
					fmt.Printf("GetMachineStageById failed: " + err.Error())
					db.Close()
					continue
				}
				if !strings.EqualFold(stage, "todo") {
					db.Close()
					continue
				}

				tp, err = database.GetTaskTypeByUUID(db, uuid)
				if err != nil {
					fmt.Printf("GetTaskTypeByUUID failed: " + err.Error())
					db.Close()
					continue
				}
			}
			db.Close()

			switch tp {
			case "schema":
				fmt.Printf("do schema logic\n")
				machine.StateMachineSchema(cluster, user, cfg, machine.ToDo, ip, uuid, 0)
			case "full":
				fmt.Printf("do full logic\n")
				machine.StateMachineSchema(cluster, user, cfg, machine.ToDo, ip, uuid, 1)
			case "all":
				fmt.Printf("do all logic\n")
			default:
				fmt.Printf("type is error\n")
			}
		}
	}()

	for {
		select {
		case <- quit:
			fmt.Printf("cancel goroutine by channel")
			return
		}
	}
}