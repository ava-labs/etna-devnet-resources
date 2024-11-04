package lib

import (
	"fmt"
	"net"
)

func FindFreePort(startFrom int) (int, error) {
	for port := startFrom; port < 65535; port++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return 0, err
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return port, nil
		}
		conn.Close()
	}
	return 0, fmt.Errorf("no free port found")
}

func FindMultipleFreePorts(count int, startFrom int) ([]int, error) {
	ports := []int{}
	for i := 0; i < count; i++ {
		port, err := FindFreePort(startFrom + i)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	return ports, nil
}
