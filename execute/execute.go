package execute

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCommand(path, cmd string, args ...string) (string, error) {
	name := path + "/" + cmd
	fmt.Printf("== command path: %s, %v\n", name, args)
	exec := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	exec.Stdout = &out
	exec.Stderr = &stderr
	err := exec.Run()
	if err != nil {
		fmt.Printf("ExecuteCommand: %s, %s, %s\n", err.Error(), out.String(), stderr.String())
		return stderr.String(), err
	}
	return stderr.String(), nil
}
