package task

import (
	"sync"
	"time"
	"fmt"
)

func Task(quit <-chan time.Time, wg *sync.WaitGroup, rate int) {
	defer wg.Done()

	checkTick := time.NewTicker(time.Duration(rate) * time.Millisecond)
	defer checkTick.Stop()
	for {
		select {
		case <- quit:
			fmt.Printf("cancel goroutine by channel")
			return
		case <- checkTick.C:
			//获取任务类型和任务状态，设置状态各个部分用 协程
			tp := "schema"
			switch tp {
			case "schema":
				fmt.Printf("do schema logic\n")

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