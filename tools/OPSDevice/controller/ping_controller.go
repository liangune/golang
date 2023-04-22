package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"go/gopkg/logger/vglog"
	"go/gopkg/middleware"
	"go/tools/OPSDevice/device"
	"go/tools/OPSDevice/ping"
	"sort"
)

const (
	DeviceSourceIsFile = "file"
	DeviceSourceIsDB   = "db"
)

func PingExecute(ctx *gin.Context) {
	cmdSlice := make(ping.PingCmdSort, 0)

	pingType, _ := ctx.GetQuery("type")
	var deviceMap map[string]device.Device
	if pingType == DeviceSourceIsFile {
		fileName := viper.GetString("deviceSource.fileName")
		sheet := viper.GetString("deviceSource.sheet")
		deviceMap = ReadExcelFile(fileName, sheet)

	} else {
		deviceMap = device.DefaultDeviceManager.PingMap
	}

	for _, v := range deviceMap {
		cmd := ping.PingCmd{
			DeviceId:   v.DeviceId,
			DeviceName: v.DeviceName,
			IpAddr:     v.IpAddr,
			Length:     56,
			Count:      4,
		}
		cmd.Ping()
		cmdSlice = append(cmdSlice, cmd)
	}
	sort.Sort(cmdSlice)

	// 创建excel表
	file := excelize.NewFile()
	sheet := "设备列表"
	// 创建一个工作表
	index := file.NewSheet(sheet)
	// 设置单元格的值
	file.SetCellValue(sheet, "A1", "设备ID")
	file.SetCellValue(sheet, "B1", "设备名称")
	file.SetCellValue(sheet, "C1", "设备IP地址")
	file.SetCellValue(sheet, "D1", "ping丢包率")
	row := 2
	for _, v := range cmdSlice {
		axisA := fmt.Sprintf("A%d", row)
		axisB := fmt.Sprintf("B%d", row)
		axisC := fmt.Sprintf("C%d", row)
		axisD := fmt.Sprintf("D%d", row)
		file.SetCellValue(sheet, axisA, v.DeviceId)
		file.SetCellValue(sheet, axisB, v.DeviceName)
		file.SetCellValue(sheet, axisC, v.IpAddr)
		file.SetCellValue(sheet, axisD, fmt.Sprintf("%d%%", v.PacketLoss))
		row++
	}
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := file.SaveAs("ping_result.xlsx"); err != nil {
		vglog.Error("%v", err)
	}

	middleware.ResponseJSON(ctx, 0, "ping命令执行完成, 结果请查看报表", nil)
}

func ReadExcelFile(fileName string, sheet string) map[string]device.Device {
	deviceMap := map[string]device.Device{}
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		vglog.Error("%v", err)
		return nil
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		vglog.Error("%v", err)
		return nil
	}
	for index, row := range rows {
		if index < 1 {
			continue
		}
		if len(row) >= 3 {
			d := device.Device{
				DeviceId:   row[0],
				DeviceName: row[1],
				IpAddr:     row[2],
				IsOnline:   device.IsOffline,
			}
			deviceMap[d.DeviceId] = d
		}
	}
	return deviceMap
}
