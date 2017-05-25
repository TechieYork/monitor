package client

import (
	"errors"
	"syscall"
)

//Global variables setting
var (
	UnixClient = NewUnixMonitorClient("/var/tmp/monitor.sock")
	UdpClient = NewUdpMonitorClient([4]byte{127, 0, 0, 1}, 5656)

	monitorSendBufferLen = 4 * 1024 * 1024
	monitorPointMaxSize = 4096
)

//Unix monitor client
type UnixMonitorClient struct {
	path string					//Unix socket path, default is /var/tmp/monitor.sock
	addr syscall.SockaddrUnix	//Socket address
	socket int					//Socket fd
	isInit bool					//Is initialized

	bufferSize int				//Socket buffer size
	pointMaxSize int 			//Point max size
}

//New function
func NewUnixMonitorClient (path string) *UnixMonitorClient {
	return &UnixMonitorClient{
		path: path,
		addr: syscall.SockaddrUnix{Name:path},
		socket: 0,
		isInit: false,
		bufferSize: monitorSendBufferLen,
		pointMaxSize: monitorPointMaxSize,
	}
}

//Init monitor unix socket
func (client *UnixMonitorClient) Init() error {
	var err error

	//Check is init or not
	if client.isInit{
		return nil
	}

	//Init socket and address
	client.socket, err = syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	//Set socket write buffer
	err = syscall.SetsockoptInt(client.socket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, client.bufferSize)

	if err != nil {
		return err
	}

	client.isInit = true

	return nil
}

//Uninit monitor unix socket
func (client *UnixMonitorClient) Uninit() error {
	if !client.isInit {
		return nil
	}

	err := syscall.Close(client.socket)

	if err != nil {
		return err
	}

	client.socket = 0
	client.isInit = false

	return nil
}

//Send monitor point via unix socket
//point string must less than 4096
func (client *UnixMonitorClient) Send(point string) error {
	if !client.isInit {
		return errors.New("Monitor not init!")
	}

	if len(point) > client.pointMaxSize {
		return errors.New("Point string size too large!")
	}

	err := syscall.Sendto(client.socket, []byte(point),0, &client.addr)

	if err != nil {
		return err
	}

	return nil
}

//Set address
func (client *UnixMonitorClient) SetAddr(path string) {
	client.path = path
	client.addr = syscall.SockaddrUnix{Name:path}
}

//Udp monitor client
type UdpMonitorClient struct {
	ip [4]byte					//Udp ip, default is [127, 0, 0, 1]
	port int					//Udp port, default is 5656
	addr syscall.SockaddrInet4	//Socket address
	socket int					//Socket fd
	isInit bool					//Is initialized

	bufferSize int				//Socket buffer size
	pointMaxSize int 			//Point max size
}

//New function
func NewUdpMonitorClient (ip [4]byte, port int) *UdpMonitorClient {
	return &UdpMonitorClient{
		ip: ip,
		port: port,
		addr: syscall.SockaddrInet4{Addr:ip, Port:port},
		socket: 0,
		isInit: false,
		bufferSize: monitorSendBufferLen,
		pointMaxSize: monitorPointMaxSize,
	}
}

//Init monitor udp socket
func (client *UdpMonitorClient) Init() error {
	var err error

	//Check is init or not
	if client.isInit {
		return nil
	}

	//Init socket and address
	client.socket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	//Set socket write buffer
	err = syscall.SetsockoptInt(client.socket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, client.bufferSize)

	if err != nil {
		return err
	}

	client.isInit = true

	return nil
}

//Uninit monitor udp socket
func (client *UdpMonitorClient) Uninit() error {
	if !client.isInit {
		return nil
	}

	err := syscall.Close(client.socket)

	if err != nil {
		return err
	}

	client.socket = 0
	client.isInit = false

	return nil
}

//Send monitor point via udp socket
//point string must less than 4096
func (client *UdpMonitorClient) Send(point string) error {
	if !client.isInit {
		return errors.New("Monitor not init!")
	}

		if len(point) > client.pointMaxSize {
		return errors.New("Point string size too large!")
	}

	err := syscall.Sendto(client.socket, []byte(point),0, &client.addr)

	if err != nil {
		return err
	}

	return nil
}

//Set address
func (client *UdpMonitorClient) SetAddr(ip [4]byte, port int) {
	client.ip = ip
	client.port = port
	client.addr = syscall.SockaddrInet4{Addr:ip, Port:port}
}
