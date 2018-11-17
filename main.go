package main

import (
	"os"

	"github.com/liuyh73/go-web-form/service"
	flag "github.com/spf13/pflag"
)

// 定义默认端口
const (
	PORT string = "8080"
)

func main() {
	// 获取用户自定义端口号
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}
	// 创建应用服务
	server := service.NewServer()
	// 运行该服务，监听port端口
	server.Run(":" + port)
}
