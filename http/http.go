package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SetBinglogEnable(url, env, clustername, username string, enable bool) error {
	body := BinLogEnableInfo{
		baseInfo{env, clustername, username},
		enable,
	}

	str, err := json.Marshal(body)
	fmt.Printf("request: %s\n", str)
	if err != nil {
		return err
	}
	res, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	if res.StatusCode != 200 {
		return errors.New("post return code is not 200")
	}
	s, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", s)
	response := ResponseInfo{}
	json.Unmarshal(s, &response)
	if response.Code != 0 {
		return errors.New("result code is not 0")
	}
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
	res, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	if res.StatusCode != 200 {
		return errors.New("post return code is not 200")
	}
	s, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", s)
	response := ResponseInfo{}
	json.Unmarshal(s, &response)
	if response.Code != 0 {
		return errors.New("result code is not 0")
	}
	fmt.Printf("res: %v\n", response)
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
	res, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	if res.StatusCode != 200 {
		return errors.New("post return code is not 200")
	}
	s, _ := ioutil.ReadAll(res.Body)
	response := ResponseInfo{}
	json.Unmarshal(s, &response)
	if response.Code != 0 {
		return errors.New("result code is not 0")
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
	res, err1 := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(str))
	if err1 != nil {
		return err1
	}
	if res.StatusCode != 200 {
		return errors.New("post return code is not 200")
	}
	s, _ := ioutil.ReadAll(res.Body)
	response := ResponseInfo{}
	json.Unmarshal(s, &response)
	if response.Code != 0 {
		return errors.New("result code is not 0")
	}
	return nil
}