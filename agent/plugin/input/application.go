package main

import(
	"fmt"
	"time"
	"net"
	"sync"
	"strings"
	"strconv"
	"errors"

	"github.com/DarkMetrix/monitor/agent/src/config"
	"github.com/DarkMetrix/monitor/agent/src/protocol"
)

var GlobalNodeInfo config.NodeInfo
var GlobalConfig map[string]string

var GlobalPointMap map[string]int

var GlobalMutex sync.Mutex

func Init(nodeInfo config.NodeInfo, config map[string]string) error {
	GlobalConfig = make(map[string]string)
	GlobalConfig = config
	GlobalNodeInfo = nodeInfo

	GlobalPointMap = make(map[string]int)

	_, ok := config["udp_address"]

	if !ok {
		return errors.New("Missing config 'udp_address'")
	}

	udpAddress, err := net.ResolveUDPAddr("udp4", config["udp_address"])

	if err != nil {
		return errors.New("Resolve udp addr failed! udp address:" + config["udp_address"])
	}

	udpConn, err := net.ListenUDP("udp4", udpAddress)

	if err != nil {
		return errors.New("Listen udp addr failed! error:" + err.Error())
	}

	err = udpConn.SetReadBuffer(16 * 1024 * 1024)

	if err != nil {
		return errors.New("Set read buffer 16M failed! error:" + err.Error())
	}

	count := 0

	go func (conn *net.UDPConn) {
		data := make([]byte, 4096)

		for {
			read, _, err := conn.ReadFromUDP(data)

			count++

			if count % 10000 == 0 {
				fmt.Println("count:", count)
			}

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
