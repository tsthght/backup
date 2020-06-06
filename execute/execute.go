package execute

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCommand(path, cmd string, args ...string) (string, error) {
	name := path + "/" + cmd
	fmt.Printf("== command path: %s\n", name)
	exec := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	exec.Stdout = &out
	exec.Stderr = &stderr
	err := exec.Run()
	if err != nil {
		return "", err
	}
	return stderr.String(), nil
}
