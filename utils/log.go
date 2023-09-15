// Copyright 2023. GodQ. All rights reserved.

package utils

import (
	"fmt"

	"github.com/go-logr/logr"
)

var (
	logger logr.Logger
)

// SetLogger sets a concrete logging implementation for all deferred Loggers.
func SetLogger(l logr.Logger) {
	logger = l
}

// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md#what-method-to-use
const (
	ERROR = 0
	WARN  = 1
	INFO  = 2
	DEBUG = 4
	TRACE = 5
)

func Log(msg string, level int) {
	logger.V(level).Info(msg)
}

func LogError(a ...any) {
	msg := fmt.Sprint(a...)
	Log(msg, ERROR)
}

func LogWarn(a ...any) {
	msg := fmt.Sprint(a...)
	Log(msg, WARN)
}

func LogInfo(a ...any) {
	msg := fmt.Sprint(a...)
	Log(msg, INFO)
}

func LogDebug(a ...any) {
	msg := fmt.Sprint(a...)
	Log(msg, DEBUG)
}

func LogTrace(a ...any) {
	msg := fmt.Sprint(a...)
	Log(msg, TRACE)
}

func LogErrorf(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	Log(msg, ERROR)
}

func LogWarnf(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	Log(msg, WARN)
}

func LogInfof(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	Log(msg, INFO)
}

func LogDebugf(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	Log(msg, DEBUG)
}

func LogTracef(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	Log(msg, TRACE)
}
