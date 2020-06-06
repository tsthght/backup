package execute

import (
	"errors"
	"fmt"
	"os/exec"
)

func ExecuteCommand(path, cmd string, args ...string) (string, error) {
	name := path + "/" + cmd
	fmt.Printf("== command path: %s\n", name)
	exec := exec.Command(name, args...)
	out, err := exec.Output()
	if err != nil {
		return "", errors.New("call output failed. err: " + err.Error())
	}
	exec.Run()
	return string(out), nil
}
