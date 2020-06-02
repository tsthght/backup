package utils

import (
	"net"
	"errors"
)

func GetLocalIP() (error, string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return err, ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return nil, ipnet.IP.String()
			}
		}
	}
	return errors.New("can not get local ip"), ""
}
