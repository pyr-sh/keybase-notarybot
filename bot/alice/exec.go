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

// Runs the command and decodes its output.
func (r *Result) DecodeOnce(input interface{}) error {
	output, err := r.RawOnce()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(output, input); err != nil {
		return err
	}
	return nil
}

// Runs the command simply returning its output.
func (r *Result) RawOnce() ([]byte, error) {
	res, err := r.cmd.Output()
	if err != nil {
		if x, ok := err.(*exec.ExitError); ok {
			return nil, errors.Wrapf(err, "stderr: %s", string(x.Stderr))
		}
		return nil, err
	}
	return res, nil
}

// Runs the command, discarding its output.
func (r *Result) RunOnce() error {
	if err := r.cmd.Run(); err != nil {
		if x, ok := err.(*exec.ExitError); ok {
			return errors.Wrapf(err, "stderr: %s", string(x.Stderr))
		}
		return err
	}
	return nil
}

// Runs the command returning its output as a reader
func (r *Result) RawStream() (io.ReadCloser, error) {
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := r.cmd.Start(); err != nil {
		return nil, err
	}
	return stdout, nil
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

// Transforms the result into a stream, monitoring both its stdout and stderr.
// Required to interact with the listen APIs.
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

// Completes the streaming process, printing out stderr if it exited with an error.
func (s *StreamedResult) Close() error {
	if err := s.eg.Wait(); err != nil {
		return errors.Wrapf(err, "stderr: %s", string(s.stderrTail))
	}
	return nil
}

// Loads up the next message into StreamedResult's internal buffer
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

// If any read error occured during processing, it's returned here
func (s *StreamedResult) Err() error { return s.readError }

// Decodes the current buffer's contents into input
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

type jm map[string]interface{}

// Executes a single Keybase CLI command without any stdin
func (c *Client) Exec(ctx context.Context, args ...interface{}) (*Result, error) {
	args = append([]interface{}{"--no-auto-fork"}, args...)
	return c.ExecWithInput(ctx, nil, args...)
}

// Executes a Keybase CLI command with the passed stdin
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
