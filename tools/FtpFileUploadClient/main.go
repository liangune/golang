package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go/gopkg/logger/vglog"
	"go/tools/FtpFileUploadClient/filehandle"
	"go/tools/FtpFileUploadClient/filescanner"
	"go/tools/FtpFileUploadClient/utils"
	"time"
)

func main() {
	//config init
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper read config failed: ", err)
		return
	}

	serverRunmode := viper.GetString("server.runmode")
	maxAge := viper.GetUint32("log.maxAge")
	// 初始化日记库
	vglog.VglogInit("./logs", vglog.InfoLog, serverRunmode)
	vglog.NewGlogCleaner(vglog.InitOption{
		Path:           "./logs/",
		Interval:       time.Minute * 30,
		Reserve:        maxAge,
		Compress:       true,
		CompressMethod: vglog.CompressMethodZip,
	})

	maxWorker := viper.GetInt("ftp.readThread")
	vglog.Info("ftp read file thread: %d", maxWorker)
	filehandle.GFileHandler = filehandle.NewFileHandler(maxWorker, filehandle.DefaultMaxTaskCount)
	filehandle.GFileHandler.SetNewHandlerFunc(filehandle.NewFileUploader)
	filehandle.GFileHandler.Start()
	vglog.Info("file handler start")

	srcPath := viper.GetString("ftp.ftpPath")
	destPath := viper.GetString("ftp.tmpPath")
	fileScanner := filescanner.NewFileScanner(srcPath, destPath)
	fileScanner.SetCallBackFunc(utils.ScanDirectoryWalkFunc)
	fileScanner.Start()
	vglog.Info("ftp src: %s, ftp tmp: %s", srcPath, destPath)
	vglog.Info(">>>>>>>>>> FtpFileUploadClient start <<<<<<<<<<")
	for {
		time.Sleep(time.Second)
	}
}
