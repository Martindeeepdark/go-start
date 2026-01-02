package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"{{.ModulePath}}/internal/application"
	"{{.ModulePath}}/internal/routes"
	"github.com/gin-gonic/gin"
)

// @title           {{.ModuleName}} API
// @version         1.0
// @description     {{.Description}}
func main() {
	ctx := context.Background()

	// Please do not change the order of the function calls below
	setCrashOutput()

	if err := loadEnv(); err != nil {
		panic("loadEnv failed, err=" + err.Error())
	}

	setLogLevel()

	if err := application.Init(ctx); err != nil {
		panic("InitializeInfra failed, err=" + err.Error())
	}

	startHttpServer()
}

// setCrashOutput 设置崩溃输出
func setCrashOutput() {
	crashFile, err := os.Create("crash.log")
	if err != nil {
		log.Printf("创建崩溃日志文件失败: %v", err)
		return
	}
	debug.SetCrashOutput(crashFile, debug.CrashOptions{})
}

// loadEnv 加载环境变量
func loadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	fileName := ".env"
	if appEnv != "" {
		fileName = ".env." + appEnv
	}

	log.Printf("加载环境变量文件: %s", fileName)
	return nil
}

// setLogLevel 设置日志级别
func setLogLevel() {
	level := getEnv("LOG_LEVEL", "info")
	log.Printf("日志级别: %s", level)
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// startHttpServer 启动 HTTP 服务器
func startHttpServer() {
	addr := getEnv("SERVER_ADDR", ":8080")
	readTimeout := getDurationEnv("READ_TIMEOUT", 15*time.Second)
	writeTimeout := getDurationEnv("WRITE_TIMEOUT", 15*time.Second)
	idleTimeout := getDurationEnv("IDLE_TIMEOUT", 60*time.Second)

	log.Printf("启动 HTTP 服务器: %s", addr)

	// 创建 Gin engine
	r := gin.Default()

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic("启动 HTTP 服务器失败: " + err.Error())
	}
}

// getDurationEnv 获取时间间隔类型的环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	return defaultValue
}
