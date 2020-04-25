package alice

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type FS struct {
	c *Client
}

var ErrNotExist = errors.New("file or folder does not exist")

var whitespaceRe = regexp.MustCompile(`\s+`)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *fileInfo) ModTime() time.Time { return fi.modTime }
func (fi *fileInfo) IsDir() bool        { return fi.mode.IsDir() }
func (fi *fileInfo) Sys() interface{}   { return nil }

type ListOpts struct {
	Recursive bool
	All       bool
}

func (c FS) List(ctx context.Context, path string, opts *ListOpts) ([]os.FileInfo, error) {
	args := []interface{}{
		"fs", "ls", path, "--long",
	}
	if opts != nil {
		if opts.All {
			args = append(args, "--all")
		}
		if opts.Recursive {
			args = append(args, "--rec")
		}
	}
	res, err := c.c.Exec(ctx, args...)
	if err != nil {
		return nil, err
	}
	output, err := res.RawOnce()
	if err != nil {
		if strings.Contains(err.Error(), "file or folder does not exist (code 5103)") {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := []os.FileInfo{}
	for i, line := range bytes.Split(output, []byte("\n")) {
		parts := whitespaceRe.Split(string(line), 6)
		if len(parts) < 6 {
			continue
		}

		mode := parseMask(parts[0])
		size, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, errors.Errorf("failed to parse size %s", parts[1])
		}

		var (
			modTime time.Time
			timeStr = parts[2] + " " + parts[3] + " " + parts[4]
		)
		if t, err := time.Parse("Jan 02 15:04", timeStr); err == nil {
			modTime = t
		}
		if t, err := time.Parse("Jan 02 2006", timeStr); err == nil {
			modTime = t
		}

		// name can be the path we're looking for, so detect that
		name := filepath.Join(path, parts[5])
		if i == 0 && filepath.Base(parts[5]) == filepath.Base(path) {
			name = path
		}

		result = append(result, &fileInfo{
			name:    name,
			size:    size,
			mode:    mode,
			modTime: modTime,
		})
	}
	return result, nil
}

func (c FS) Mkdir(ctx context.Context, path string) error {
	res, err := c.c.Exec(ctx, "fs", "mkdir", path)
	if err != nil {
		return err
	}
	return res.RunOnce()
}

type WriteOpts struct {
	Append  bool
	BufSize int
}

func (c FS) Write(ctx context.Context, path string, input io.Reader, opts *WriteOpts) error {
	args := []interface{}{"fs", "write", path}
	if opts != nil {
		if opts.Append {
			args = append(args, "--append")
		}
		if opts.BufSize > 0 {
			args = append(args, "--buffersize", strconv.Itoa(opts.BufSize))
		}
	}

	res, err := c.c.ExecWithInput(ctx, input, args...)
	if err != nil {
		return err
	}
	if err := res.RunOnce(); err != nil {
		return err
	}
	return nil
}
