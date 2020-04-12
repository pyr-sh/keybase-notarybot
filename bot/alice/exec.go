package alice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

type Result struct {
	cmd *exec.Cmd
}

func (r *Result) DecodeOnce(input interface{}) error {
	output, err := r.cmd.Output()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(output, input); err != nil {
		return err
	}
	return nil
}

func (r *Result) RunOnce() error {
	/*
		a, err := r.cmd.CombinedOutput()
		if err != nil {
			return err
		}
		return nil
	*/
	return r.cmd.Run()
}

func (c *Client) commonArgs() []string {
	args := []string{}
	if c.cfg.homePath != "" {
		args = append(args, "-H", c.cfg.homePath)
	}
	return args
}

type J map[string]interface{}

func (c *Client) Exec(ctx context.Context, args ...interface{}) (*Result, error) {
	return c.ExecWithInput(ctx, nil, args...)
}

func (c *Client) ExecWithInput(ctx context.Context, body io.Reader, args ...interface{}) (*Result, error) {
	commandArgs := c.commonArgs()
	for _, arg := range args {
		commandArgs = append(commandArgs, fmt.Sprintf("%s", arg))
	}

	cmd := exec.CommandContext(ctx, c.cfg.executablePath, commandArgs...)
	if body != nil {
		cmd.Stdin = body
	}
	res := &Result{
		cmd: cmd,
	}
	return res, nil
}
