package openports

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
		// Try to listen on the port instead of dialing it
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			continue // Port is in use, try next one
		}
		listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no free port found")
}

func FindMultipleFreePorts(count int, startFrom int) ([]int, error) {
	ports := []int{}
	lastPort := startFrom - 1
	for i := 0; i < count; i++ {
		port, err := FindFreePort(lastPort + 1)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
		lastPort = port
	}
	return ports, nil
}
