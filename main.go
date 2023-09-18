package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/godq/go-rest-agent/agent"
	"github.com/godq/go-rest-agent/utils"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	port := "5000"
	args := os.Args
	if len(args) == 1 {
		port = "5000"
	} else if len(args) == 2 {
		port = args[1]
	} else {
		fmt.Println("At most 1 arg (port) is allowed.")
	}

	opts := zap.Options{
		Development: true,
		DestWriter:  os.Stdout,
	}
	logger := zap.New(zap.UseFlagOptions(&opts))
	utils.SetLogger(logger)
	r := gin.Default()
	agent.SetApi(r)
	r.Run("0.0.0.0:" + port)
}
