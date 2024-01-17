package main

import (
	"log"
	"net/http"
	"time"

	"configuration-management/global"
	"configuration-management/internal/routers"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"
	"configuration-management/pkg/setting"
	"configuration-management/utils/security"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title 用户配置管理中心
// @version 1.0
// @description 统一管理用户的配置信息
func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()

	server := http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		// 打印日志
		global.Logger.Panicf("server.ListenAndServe err: %v", err, errcode.ServerError.WithDetails(err.Error()))
	}
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupRSAKey()
	if err != nil {
		log.Fatalf("init.setupRSAKey err: %v", err)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupDBEngine() error {
	var err error
	dsn := global.DatabaseSetting.UserName + ":" + global.DatabaseSetting.Password + "@tcp(" + global.DatabaseSetting.Host + ")/" + global.DatabaseSetting.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	global.DBEngine, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func setupRSAKey() error {
	var err error
	global.PrivateKey, err = security.LoadPrivateKey(global.AppSetting.PrivateKeyFilePath)
	if err != nil {
		return err
	}
	global.PublicKey, err = security.LoadPublicKey(global.AppSetting.PublicKeyFilePath)
	if err != nil {
		return err
	}

	return nil
}
