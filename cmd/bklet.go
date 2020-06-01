package main

import (
	"fmt"

	"github.com/tsthght/backup/http"
)

func main() {
	fmt.Printf("i am bklet")
	http.SetBinglogEnable(
		"http://xxxxxx:8000/api/v1/cluster/conf_cluster_binlog",
		"product",
		"inf_blade_multiidc",
		"xxxx",
		true,
		)
	/*
	wg := sync.WaitGroup{}
	wg.Add(1)
	quit := make(chan time.Time)
	go register.Register(quit, &wg, 1000)
	wg.Wait()

	 */
}
