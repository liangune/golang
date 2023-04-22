package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go/gopkg/db"
	"go/gopkg/logger/vglog"
	"go/tools/OPSDevice/device"
	"go/tools/OPSDevice/global"
	"go/tools/OPSDevice/kafkaconsumer"
	"go/tools/OPSDevice/router"
	"strconv"
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
	}

	serverPort := viper.GetInt64("server.port")
	serverPortStr := strconv.Itoa(int(serverPort))

	serverRunmode := viper.GetString("server.runmode")
	// 初始化日记库
	vglog.VglogInit("./logs", vglog.InfoLog, serverRunmode)
	vglog.NewGlogCleaner(vglog.InitOption{
		Path:           "./logs/",
		Interval:       time.Minute * 30,
		Reserve:        0,
		Compress:       true,
		CompressMethod: vglog.CompressMethodZip,
	})

	// 初始化数据库
	global.DbPool, err = db.DBPoolInit(&db.DBConfig{
		Username: viper.GetString("postgresql.username"),
		Password: viper.GetString("postgresql.password"),
		Host:     viper.GetString("postgresql.host"),
		Port:     viper.GetInt("postgresql.port"),
		Dbname:   viper.GetString("postgresql.dbname"),
		DbType:   db.PostgreSQL,
	})
	if err != nil {
		fmt.Println("db pool init is error: ", err)
		return
	}
	defer global.DbPool.Close()

	err = device.DefaultDeviceManager.Init()
	if err != nil {
		fmt.Println("get device error: ", err)
	}

	topic := viper.GetString("kafka.topic")
	brokers := viper.GetString("kafka.brokers")
	kafkaconsumer.DefaultConsumer = kafkaconsumer.NewConsumer(brokers, topic, "groupIdOpsDevice")

	_ = NewReportManager()

	gin.SetMode(viper.GetString("server.runmode"))
	engine := gin.New()
	router.SetupRouter(engine)
	engine.Run(":" + serverPortStr)
}
