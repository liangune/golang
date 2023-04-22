package ping

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go/gopkg/logger/vglog"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func PingExecCmd(destIP string, count int) {
	sCount := strconv.Itoa(count)
	var buf bytes.Buffer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", destIP)
	} else {
		cmd = exec.Command("ping", "-c", sCount, destIP)
	}

	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		vglog.Error("ping %s -c %s fail: %v", destIP, sCount, err)
	} else {
		vglog.Info("%s", buf.String())
	}
}

const (
	MAX_PING_SEND_LEN    = 2000
	MAX_PING_RECEIVE_LEN = 1000
	MAX_PING_TIMEOUT     = 3
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

type PingCmd struct {
	DeviceId   string
	DeviceName string
	IpAddr     string
	Length     int
	Count      int
	PacketLoss int
	MaxTime    float64
	MinTime    float64
	AvgTime    float64
}

func (cmd *PingCmd) CheckSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index]) << 8
	}
	rt = uint16(sum) + uint16(sum>>16)

	return ^rt
}

func (cmd *PingCmd) Ping() {
	ipAddr := cmd.IpAddr
	length := cmd.Length
	count := cmd.Count
	originBytes := make([]byte, MAX_PING_SEND_LEN)
	var (
		icmp                      ICMP
		laddr                     = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		raddr, _                  = net.ResolveIPAddr("ip", ipAddr)
		maxTime, minTime, avgTime float64
	)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()
	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	binary.Write(&buffer, binary.BigEndian, originBytes[0:length])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], cmd.CheckSum(b))
	fmt.Printf("\n正在 Ping %s 具有 %d(%d) 字节的数据:\n", raddr.String(), length, length+28)
	recv := make([]byte, MAX_PING_RECEIVE_LEN)
	retList := []float64{}
	packetLoss := 0
	maxTime = 3000.0
	minTime = 0.0
	avgTime = 0.0
	for i := 0; i < count; i++ {
		if _, err := conn.Write(buffer.Bytes()); err != nil {
			packetLoss++
			time.Sleep(time.Second)
			continue
		}
		startTime := time.Now()
		conn.SetReadDeadline((time.Now().Add(time.Second * MAX_PING_TIMEOUT)))
		_, err := conn.Read(recv)
		if err != nil {
			packetLoss++
			time.Sleep(time.Second)
			continue
		}
		endTime := time.Now()
		dur := float64(endTime.Sub(startTime).Nanoseconds()) / 1e6
		retList = append(retList, dur)
		if dur < maxTime {
			maxTime = dur
		}
		if dur > minTime {
			minTime = dur
		}
		fmt.Printf("from %s ack: icmp_seq=%d time= %.3fms\n", raddr.String(), i+1, dur)
		time.Sleep(time.Second)
	}
	cmd.PacketLoss = (packetLoss / count) * 100
	fmt.Printf("丢包率: %d%%\n", cmd.PacketLoss)
	if len(retList) == 0 {
		avgTime = 3000.0
	} else {
		sum := 0.0
		for _, n := range retList {
			sum += n
		}
		avgTime = sum / float64(len(retList))
	}
	fmt.Printf("rtt min/avg/max = %.3fms/%.3fms/%.3fms\n", minTime, avgTime, maxTime)
	cmd.MaxTime = maxTime
	cmd.MinTime = minTime
	cmd.AvgTime = avgTime
}
