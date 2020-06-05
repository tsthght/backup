package execute

import (
	"errors"
	"os/exec"
)

func ExecuteCommand(path, cmd string, args ...string) (string, error) {
	name := path + "/" + cmd
	exec := exec.Command(name, args...)
	out, err := exec.Output()
	if err != nil {
		return "", errors.New("call output failed. err: " + err.Error())
	}
	exec.Run()
	return string(out), nil
}
