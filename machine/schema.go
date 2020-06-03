package machine

import "time"

const (
	todo = iota
	prepareenv
	precheck
	dumping
	loading
	poscheck
	resetenv
)

func StateMachineSchema() int {
	for {
		time.Sleep(10000)
	}
}
