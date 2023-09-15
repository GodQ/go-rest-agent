// Copyright 2023. GodQ. All rights reserved.
package utils

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"time"
)

type CmdResult struct {
	EventChan     chan string
	OutputBuilder strings.Builder
	Code          int
}

func NewCmdResult() *CmdResult {
	cr := &CmdResult{}
	cr.EventChan = make(chan string, 1024*10)
	cr.OutputBuilder = strings.Builder{}
	return cr
}

func (cr *CmdResult) Output() string {
	return cr.OutputBuilder.String()
}

func (cr *CmdResult) Success() bool {
	if cr.Code == 0 {
		return true
	} else {
		return false
	}
}

func RunShellWithRetry(script string, result *CmdResult, timeoutSeconds int, maxRetries int) (*CmdResult, error) {
	args := []string{
		"-c",
		script,
	}
	return RunCmdWithRetry("bash", args, result, timeoutSeconds, maxRetries)
}

func RunCmdWithRetry(executable string, args []string, result *CmdResult, timeoutSeconds int, maxRetries int) (*CmdResult, error) {
	var err error
	if result == nil {
		result = NewCmdResult()
	}
	for i := 0; i <= maxRetries; i++ {
		// fmt.Println(executable, args)
		err = RunCmd(executable, args, result, timeoutSeconds)
		if err != nil {
			LogWarn(err)
			if _, ok := err.(*exec.ExitError); !ok {
				return result, err
			}
			LogWarn("Retrying")
			continue
		} else {
			return result, nil
		}
	}
	return result, err
}

// RunCli
func RunCmd(executable string, args []string, result *CmdResult, timeoutSeconds int) error {
	if result == nil {
		result = NewCmdResult()
	}
	var cmdContext context.Context
	if timeoutSeconds > 0 {
		var cancel context.CancelFunc
		cmdContext, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
		defer cancel()
	} else {
		cmdContext = context.Background()
	}

	cmd := exec.CommandContext(cmdContext, executable, args...)
	LogDebug("path: ", cmd.Path)
	LogDebug("args: ", cmd.Args)

	LogDebug("Start to run cmd")

	var cmdErr error

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.EventChan <- err.Error()
		LogError(err)
		return err
	}
	// cmd.Stderr = cmd.Stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		result.EventChan <- err.Error()
		LogError(err)
		return err
	}

	if err := cmd.Start(); err != nil {
		result.EventChan <- err.Error()
		return err
	}

	stdScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for stdScanner.Scan() {
			line := stdScanner.Text()
			result.EventChan <- line
			result.OutputBuilder.WriteString(line + "\n")
		}
	}()
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			result.EventChan <- line
			result.OutputBuilder.WriteString(line + "\n")
		}
	}()
	cmdErr = cmd.Wait()
	LogInfo("Finish running cmd")
	statusCode := cmd.ProcessState.ExitCode()
	result.Code = statusCode
	if statusCode != 0 {
		LogErrorf("cmd returned error code %d, \ndetails: %s", statusCode, cmdErr)
	}
	return cmdErr
}
