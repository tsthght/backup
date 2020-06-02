package status

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsthght/backup/database"
)

func Status(quit <-chan time.Time, wg *sync.WaitGroup, rate int, cluster *database.MGRInfo, user database.UserInfo) {
	defer wg.Done()

	checkTick := time.NewTicker(time.Duration(rate) * time.Millisecond)
	defer checkTick.Stop()
	for {
		select {
		case <- quit:
			fmt.Printf("cancel goroutine by channel")
			return
		case <- checkTick.C:

		}
	}
}
