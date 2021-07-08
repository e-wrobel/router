package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/e-wrobel/router/handlers"
	"github.com/gin-gonic/gin"
)

var (
	listen      string
	remoteURL   string
	logSeverity string
)

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:8080", "-listen 0.0.0.0:8080")
	flag.StringVar(&remoteURL, "remoteURL", "", "-remoteURL http://192.168.0.1:8888")
	flag.StringVar(&logSeverity, "logSeverity", "Info", "-logSeverity Info|Debug|Warn")
	flag.Parse()
}

func main() {

	handle, err := handlers.New(remoteURL, logSeverity)
	if err != nil {
		fmt.Printf("Error preparing parameters: %v", err)
		os.Exit(1)
	}

	service := gin.Default()
	service.Any("/*path", handle.HandleAnyRoute)

	err = service.Run(listen)
	if err != nil {
		panic("Unable to listen!")
	}
}
