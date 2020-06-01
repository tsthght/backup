package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SetBinglogEnable(url, env, clustername, username string, enable bool) error {
	body := BinLogEnableInfo{
		baseInfo{env, clustername, username},
		enable,
	}

	str, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	var r []byte
	res.Body.Read(r)
	response := ResponseInfo{}
	json.Unmarshal(r, &response)
	fmt.Printf("res: %v\n", response)
	return nil
}

func SetPumpStatus(url, env, clustername, username, command string, list []string) error {
	body := PumpInfo {
		baseInfo{env, clustername, username},
		command,
		list,
	}
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	return nil
}

func SetDrainerStatus(url, env, clustername, username, command string, list []string) error {
	body := DrainerInfo {
		baseInfo{env, clustername, username},
		command,
		list,
	}
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	return nil
}

func RollingCluster(url, env, clustername, username, command string) error {
	body := RollbackInfo {
		baseInfo{env, clustername, username},
		command,
	}
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	return nil
}