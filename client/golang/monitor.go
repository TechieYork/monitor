package client

import (
	"errors"
	"syscall"
)

//Global variables setting
var (
	monitorUnixPath = "/var/tmp/monitor.sock"
	monitorUnixAddr = syscall.SockaddrUnix{Name:monitorUnixPath}
	monitorUnixSocket = 0
	monitorUnixIsInit = false

	monitorUdpIp = [4]byte{127, 0, 0, 1}
	monitorUdpPort = 5656
	monitorUdpAddr = syscall.SockaddrInet4{Port:monitorUdpPort, Addr:monitorUdpIp}
	monitorUdpSocket = 0
	monitorUdpIsInit = false

	monitorSendBufferLen = 4 * 1024 * 1024
	monitorPointMaxSize = 4096
)

//Init monitor unix socket
func InitMonitorUnix () error {
	var err error

	//Check is init or not
	if monitorUnixIsInit{
		return nil
	}

	//Init socket and address
	monitorUnixSocket, err = syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	//Set socket write buffer
	err = syscall.SetsockoptInt(monitorUnixSocket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, monitorSendBufferLen)

	if err != nil {
		return err
	}

	monitorUnixIsInit = true

	return nil
}

//Uninit monitor unix socket
func UninitMonitorUnix () error {
	if !monitorUnixIsInit {
		return nil
	}

	err := syscall.Close(monitorUnixSocket)

	if err != nil {
		return err
	}

	monitorUnixSocket = 0

	monitorUnixIsInit = false
	return nil
}

//Send monitor point via unix socket
//point string must less than 4096
func SendMonitorUnix (point string) error {
	if !monitorUnixIsInit {
		return errors.New("Monitor not init!")
	}

	if len(point) > monitorPointMaxSize{
		return errors.New("Point string size too large!")
	}

	err := syscall.Sendto(monitorUnixSocket, []byte(point),0, &monitorUnixAddr)

	if err != nil {
		return err
	}

	return nil
}

//Init monitor udp socket
func InitMonitorUdp() error {
	var err error

	//Check is init or not
	if monitorUdpIsInit{
		return nil
	}

	//Init socket and address
	monitorUdpSocket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	//Set socket write buffer
	err = syscall.SetsockoptInt(monitorUdpSocket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, monitorSendBufferLen)

	if err != nil {
		return err
	}

	monitorUdpIsInit = true

	return nil
}

//Uninit monitor udp socket
func UninitMonitorUdp () error {
	if !monitorUdpIsInit {
		return nil
	}

	err := syscall.Close(monitorUdpSocket)

	if err != nil {
		return err
	}

	monitorUdpSocket = 0

	monitorUdpIsInit = false
	return nil
}

//Send monitor point via udp socket
//point string must less than 4096
func SendMonitorUdp (point string) error {
	if !monitorUdpIsInit {
		return errors.New("Monitor not init!")
	}

		if len(point) > monitorPointMaxSize{
		return errors.New("Point string size too large!")
	}

	err := syscall.Sendto(monitorUdpSocket, []byte(point),0, &monitorUdpAddr)

	if err != nil {
		return err
	}

	return nil
}

