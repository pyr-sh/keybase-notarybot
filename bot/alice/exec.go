package alice

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Result struct {
	ctx context.Context
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
	if err := r.cmd.Run(); err != nil {
		if x, ok := err.(*exec.ExitError); ok {
			return errors.Wrapf(err, "stderr: %s", string(x.Stderr))
		}
		return err
	}
	return nil
}

type StreamedResult struct {
	r            *Result
	Context      context.Context
	eg           *errgroup.Group
	stderrTail   []byte
	stdoutBuffer *bufio.Reader
	stdoutLine   []byte
	readError    error
}

func (r *Result) Stream() (*StreamedResult, error) {
	s := &StreamedResult{
		r: r,
	}

	stdoutOutputPipe, stdoutInputPipe := io.Pipe()
	s.stdoutBuffer = bufio.NewReader(stdoutOutputPipe)
	r.cmd.Stdout = stdoutInputPipe

	stderrOutputPipe, stderrInputPipe := io.Pipe()
	r.cmd.Stderr = stderrInputPipe

	if err := s.r.cmd.Start(); err != nil {
		return nil, err
	}

	s.eg, s.Context = errgroup.WithContext(r.ctx)
	s.eg.Go(func() error {
		defer stderrOutputPipe.Close()

		buf := make([]byte, 1024)
		for {
			select {
			case <-s.Context.Done():
				return context.Canceled
			default:
			}

			n, err := stderrOutputPipe.Read(buf)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			s.stderrTail = append(s.stderrTail, buf[:n]...)
			if len(s.stderrTail) > 500 {
				s.stderrTail = s.stderrTail[len(s.stderrTail)-500:]
			}
		}
	})
	s.eg.Go(func() error {
		defer stdoutOutputPipe.Close()

		// Main goroutine monitors the process itself, kills stdout
		if err := s.r.cmd.Wait(); err != nil {
			return err
		}
		return nil
	})

	return s, nil
}

func (s *StreamedResult) Close() error {
	if err := s.eg.Wait(); err != nil {
		return errors.Wrapf(err, "stderr: %s", string(s.stderrTail))
	}
	return nil
}

func (s *StreamedResult) Next() bool {
	line := []byte{}
	for {
		select {
		case <-s.Context.Done():
			s.readError = context.Canceled
			return false
		default:
		}

		b, isPrefix, err := s.stdoutBuffer.ReadLine()
		if err != nil {
			if err == io.EOF {
				return false
			}

			s.readError = err
			return false
		}
		line = append(line, b...)
		if !isPrefix {
			break
		}
	}
	s.stdoutLine = line
	return true
}
func (s *StreamedResult) Err() error { return s.readError }
func (s *StreamedResult) Decode(input interface{}) error {
	return json.Unmarshal(s.stdoutLine, input)
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
		if str, ok := arg.(string); ok {
			commandArgs = append(commandArgs, str)
			continue
		}
		if stringer, ok := arg.(fmt.Stringer); ok {
			commandArgs = append(commandArgs, stringer.String())
			continue
		}
		encoded, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		commandArgs = append(commandArgs, string(encoded))
	}

	cmd := exec.CommandContext(ctx, c.cfg.executablePath, commandArgs...)
	if body != nil {
		cmd.Stdin = body
	}
	res := &Result{
		ctx: ctx,
		cmd: cmd,
	}
	return res, nil
}
