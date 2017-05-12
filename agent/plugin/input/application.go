package main

import(
	"time"
	"net"
	"sync"
	"strings"
	"strconv"
	"errors"
	"os"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalPointMap map[string]int

var GlobalMutex sync.Mutex

func initUdp(addr string) error {
	udpAddress, err := net.ResolveUDPAddr("udp4", addr)

	if err != nil {
		return errors.New("Resolve udp addr failed! udp address:" + addr)
	}

	udpConn, err := net.ListenUDP("udp4", udpAddress)

	if err != nil {
		return errors.New("Listen udp addr failed! error:" + err.Error())
	}

	err = udpConn.SetReadBuffer(16 * 1024 * 1024)

	if err != nil {
		return errors.New("Set read buffer 16M failed! error:" + err.Error())
	}

	go func (conn *net.UDPConn) {
		data := make([]byte, 4096)

		for {
			read, _, err := conn.ReadFromUDP(data)

			if err != nil {
				continue
			}

			point := strings.Split(string(data[:read]), "|")

			if len(point) != 2 {
				continue
			}

			key := point[0]
			value, err := strconv.Atoi(point[1])

			if err != nil {
				continue
			}

			GlobalMutex.Lock()

			oldValue, ok := GlobalPointMap[key]

			if !ok {
				GlobalPointMap[key] = value
			} else {
				GlobalPointMap[key] = oldValue + value
			}

			GlobalMutex.Unlock()
		}

	}(udpConn)

	return nil
}

func initUnix(addr string) error {
	unixAddr, err := net.ResolveUnixAddr("unixgram", addr)

	if nil != err {
		return errors.New("Resolve unix addr failed! unix address:" + addr)
	}

	unixConn, err := net.ListenUnixgram("unixgram", unixAddr)

	if nil != err {
		return errors.New("Listen unix addr failed! error:" + err.Error())
	}

	//Change unix domain socket file mode
	err = os.Chmod(addr, 0777)

	if err != nil {
		return errors.New("Chmod on " + addr + " failed! error:" + err.Error())
	}

	//Change unix domain socket file to nobody
	err = os.Chown(addr, 99, 99)

	if err != nil {
		return errors.New("Chown on " + addr + " failed! error:" + err.Error())
	}

	go func (conn *net.UnixConn) {
		data := make([]byte, 4096)

		for {
			read, err := conn.Read(data)

			if err != nil {
				continue
			}

			point := strings.Split(string(data[:read]), "|")

			if len(point) != 2 {
				continue
			}

			key := point[0]
			value, err := strconv.Atoi(point[1])

			if err != nil {
				continue
			}

			GlobalMutex.Lock()

			oldValue, ok := GlobalPointMap[key]

			if !ok {
				GlobalPointMap[key] = value
			} else {
				GlobalPointMap[key] = oldValue + value
			}

			GlobalMutex.Unlock()
		}

	}(unixConn)

	return nil
}

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	GlobalPointMap = make(map[string]int)

	udpAddr, udpAddrOk := config["udp_address"]
	unixAddr, unixAddrOk := config["unix_address"]

	if !udpAddrOk && !unixAddrOk {
		return errors.New("Missing config 'udp_address' or 'unix_address'")
	}

	if udpAddrOk {
		err := initUdp(udpAddr)

		if err != nil {
			return err
		}
	}

	if unixAddrOk {
		err := initUnix(unixAddr)

		if err != nil {
			return err
		}
	}

	return nil
}

func Collect()(*protocol.Proto, error) {
	proto := protocol.NewProto(1)

	curTime := time.Now()
	currentTime := curTime.Local().Format("2006-01-02 15:04:05")

	GlobalMutex.Lock()
	defer GlobalMutex.Unlock()

	for pointKey, pointValue := range GlobalPointMap{
		data := protocol.NewData()
		data.Time = currentTime

		data.Tag["node_name"] = GlobalNodeInfo.Name
		data.Tag["node_ip"] = GlobalNodeInfo.IP
		data.Field[pointKey] = pointValue

		proto.DataList = append(proto.DataList, *data)
	}

	GlobalPointMap = make(map[string]int)

	return proto, nil
}
