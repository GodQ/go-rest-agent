package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/godq/go-rest-agent/agent"
	"github.com/godq/go-rest-agent/utils"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{
		Development: true,
		DestWriter:  os.Stdout,
	}
	logger := zap.New(zap.UseFlagOptions(&opts))
	utils.SetLogger(logger)
	r := gin.Default()
	agent.SetApi(r)
	r.Run("0.0.0.0:5000")
}
