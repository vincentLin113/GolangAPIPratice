package main

import (
	"fmt"
	"os"
	"syscall"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/gredis"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/setting"
	"vincent-gin-go/routers"

	"github.com/fvbock/endless"
)

func main() {
	setting.Setup()
	logging.Setup()
	models.Setup()
	gredisErr := gredis.Setup()
	if gredisErr != nil {
		panic(gredisErr)
	}
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT not set")
		port = fmt.Sprintf("%d", setting.ServerSetting.HttpPort)
	}
	endPoint := fmt.Sprintf(":%s", port)
	fmt.Printf("### END POINT: %s ###\n", endPoint)
	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		fmt.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := server.ListenAndServe()
	if err != nil {
		logging.Error(err)
		panic(err)
	}
}
