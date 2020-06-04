package machine

import "time"

const (
	ToDo = iota
	PrepareEnv
	PreCheck
	Dumping
	Loading
	PosCheck
	ResetEnv
)



func StateMachineSchema(initState int) int {
	for {
		time.Sleep(10000)
	}
}
