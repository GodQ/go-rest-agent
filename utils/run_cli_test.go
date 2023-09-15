package utils

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	exe := "bash"
	args := []string{
		"-c",
		"ls -lh",
		// "while true; do echo \"123\"; sleep 2; done",
		// "while true; do date; ls xxxaaa; sleep 2; done",
	}
	result := NewCmdResult()
	err := RunCmd(exe, args, result, 5)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	// t.Logf("%+v", result)
	t.Log(result.Output())
}

func TestRunShell(t *testing.T) {
	script := "ls -lh | base64"
	result := NewCmdResult()
	result, err := RunShellWithRetry(script, result, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("%+v", result)
	t.Log(result.Output())
}
