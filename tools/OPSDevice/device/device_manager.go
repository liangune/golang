package device

import (
	"fmt"
	"go/gopkg/logger/vglog"
	"go/tools/OPSDevice/global"
	"strconv"
	"strings"
)

type DeviceManager struct {
	// ping
	PingMap map[string]Device
	// GAT 1400
	GATMap map[string]Device
	// SDK
	SDKMap map[string]Device
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		PingMap: map[string]Device{},
		GATMap:  map[string]Device{},
		SDKMap:  map[string]Device{},
	}
}

var DefaultDeviceManager *DeviceManager = NewDeviceManager()

func (m *DeviceManager) Init() error {
	type DBDevice struct {
		Apeid               string
		Name                string
		Ipaddr              string
		Connectionagreement string
	}
	var deviceSlice []DBDevice
	sql := fmt.Sprintf("SELECT apeid, name, ipaddr, connectionagreement FROM tb_ape WHERE isonline = '1' AND accessnetworktype = 2 AND camerastatus = 1")
	err := global.DbPool.GetDB().Raw(sql).Scan(&deviceSlice).Error
	if err != nil {
		vglog.Error("get device from database: %v", err)
		return err
	}
	for _, device := range deviceSlice {
		protocolSlice := strings.Split(device.Connectionagreement, ",")
		for _, sprotocol := range protocolSlice {
			protocol, _ := strconv.Atoi(sprotocol)
			d := Device{
				DeviceId:   device.Apeid,
				DeviceName: device.Name,
				IpAddr:     device.Ipaddr,
				Protocol:   protocol,
				IsOnline:   IsOffline,
			}
			m.PingMap[d.DeviceId] = d
			if protocol == Protocol1400 {
				m.GATMap[d.DeviceId] = d
			}
			if protocol == ProtocolSDK {
				m.SDKMap[d.DeviceId] = d
			}
		}
	}
	return nil
}

func (m *DeviceManager) GetGAT1400Slice() (onlineSlice []Device, offlineSlice []Device) {
	for _, v := range m.GATMap {
		if v.IsOnline == IsOnline {
			onlineSlice = append(onlineSlice, v)
		} else {
			offlineSlice = append(offlineSlice, v)
		}
	}
	return
}

func (m *DeviceManager) GetSDKSlice() (onlineSlice []Device, offlineSlice []Device) {
	for _, v := range m.SDKMap {
		if v.IsOnline == IsOnline {
			onlineSlice = append(onlineSlice, v)
		} else {
			offlineSlice = append(offlineSlice, v)
		}
	}
	return
}
