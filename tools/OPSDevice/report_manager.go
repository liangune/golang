package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go/gopkg/logger/vglog"
	"go/tools/OPSDevice/device"
	"time"
)

const (
	defaultUpdateInterval = 10 * time.Minute
)

type ReportManager struct {
	updateInterval time.Duration
}

func NewReportManager() *ReportManager {
	m := &ReportManager{
		updateInterval: defaultUpdateInterval,
	}
	go m.updateReportTimer()
	return m
}

func (m *ReportManager) updateReportTimer() {
	ticker := time.NewTicker(m.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.CreateReportDoc(device.Protocol1400)
			m.CreateReportDoc(device.ProtocolSDK)
		}
	}
}

func (m *ReportManager) CreateReportDoc(protocol int) {
	onlineSlice := make([]device.Device, 0)
	offlineSlice := make([]device.Device, 0)
	fileName := ""
	if protocol == device.Protocol1400 {
		fileName = "1400_device.xlsx"
		onlineSlice, offlineSlice = device.DefaultDeviceManager.GetGAT1400Slice()
	} else if protocol == device.ProtocolSDK {
		fileName = "sdk_device.xlsx"
		onlineSlice, offlineSlice = device.DefaultDeviceManager.GetSDKSlice()
	}

	fmt.Println(len(onlineSlice), " CreateReportDoc ", len(offlineSlice))

	// 创建excel表
	file := excelize.NewFile()
	sheet := "设备列表"
	// 创建一个工作表
	index := file.NewSheet(sheet)
	// 设置单元格的值
	file.SetCellValue(sheet, "A1", "设备ID")
	file.SetCellValue(sheet, "B1", "设备名称")
	file.SetCellValue(sheet, "C1", "设备IP地址")
	file.SetCellValue(sheet, "D1", "是否在线")
	row := 2
	for _, v := range offlineSlice {
		axisA := fmt.Sprintf("A%d", row)
		axisB := fmt.Sprintf("B%d", row)
		axisC := fmt.Sprintf("C%d", row)
		axisD := fmt.Sprintf("D%d", row)
		file.SetCellValue(sheet, axisA, v.DeviceId)
		file.SetCellValue(sheet, axisB, v.DeviceName)
		file.SetCellValue(sheet, axisC, v.IpAddr)
		file.SetCellValue(sheet, axisD, v.IsOnline)
		row++
	}

	for _, v := range onlineSlice {
		axisA := fmt.Sprintf("A%d", row)
		axisB := fmt.Sprintf("B%d", row)
		axisC := fmt.Sprintf("C%d", row)
		axisD := fmt.Sprintf("D%d", row)
		file.SetCellValue(sheet, axisA, v.DeviceId)
		file.SetCellValue(sheet, axisB, v.DeviceName)
		file.SetCellValue(sheet, axisC, v.IpAddr)
		file.SetCellValue(sheet, axisD, v.IsOnline)
		row++
	}

	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		vglog.Error("%v", err)
	}
}
