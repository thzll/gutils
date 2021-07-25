package myutils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type InCmdWorker interface {
	Command(string, ...string) ([]byte, error)
	CommandWithContext(context.Context, string, ...string) ([]byte, error)
}

var (
	cmdWorker  CmdWorker = CmdWorker{}
	Timeout              = 3 * time.Second
	ErrTimeout           = fmt.Errorf("command timed out")
)

func Command(name string, arg ...string) ([]byte, error) {
	return cmdWorker.Command(name, arg...)
}

type CmdWorker struct {
	TimeOut time.Duration
}

func (i CmdWorker) Command(name string, arg ...string) ([]byte, error) {
	if i.TimeOut == 0 {
		i.TimeOut = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), i.TimeOut)
	defer cancel()
	return i.CommandWithContext(ctx, name, arg...)
}

func (i CmdWorker) CommandWithContext(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return buf.Bytes(), err
	}

	if err := cmd.Wait(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}
