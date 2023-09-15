package utils

import (
	"os"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestMain(m *testing.M) {
	opts := zap.Options{
		Development: true,
		DestWriter:  os.Stdout,
	}
	logger := zap.New(zap.UseFlagOptions(&opts))
	SetLogger(logger)
	m.Run()
}
