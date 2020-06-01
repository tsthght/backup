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
			fmt.Printf("do sth.\n")
		}
	}
}