package golang

import (
	"errors"
	"syscall"
)

//Global variables setting
var (
	monitorPath = "/var/tmp/monitor.sock"
	monitorAddr = syscall.SockaddrUnix{Name:monitorPath}
	monitorSocket = 0
	monitorSendBufferLen = 4 * 1024 * 1024
	monitorPointMaxSize = 4096
	monitorIsInit = false
)

//Init monitor unix socket
func InitMonitorUnix () error {
	var err error

	//Check is init or not
	if monitorIsInit{
		return nil
	}

	//Init socket and address
	monitorSocket, err = syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	//Set socket write buffer
	err = syscall.SetsockoptInt(monitorSocket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, monitorSendBufferLen)

	if err != nil {
		return err
	}

	monitorIsInit = true

	return nil
}

//Uninit monitor unix socket
func UninitMonitorUnix () error {
	if !monitorIsInit {
		return nil
	}

	err := syscall.Close(monitorSocket)

	if err != nil {
		return err
	}

	monitorSocket = 0

	monitorIsInit = false
	return nil
}

//Send monitor point via unix socket
//point string must less than 4096
func SendMonitorUnix (point string) error {
	if !monitorIsInit {
		return errors.New("Monitor not init!")
	}

	if len(point) > monitorPointMaxSize{
		return errors.New("Point string size too large!")
	}

	err := syscall.Sendto(monitorSocket, []byte(point),0, &monitorAddr)

	if err != nil {
		return err
	}

	return nil
}
