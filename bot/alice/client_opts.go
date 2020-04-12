package alice

import "github.com/mitchellh/go-homedir"

func mustExpand(input string) string {
	path, err := homedir.Expand(input)
	if err != nil {
		panic(err)
	}
	return path
}

type ClientOption interface {
	apply(*clientConfig) error
}

type clientConfig struct {
	executablePath string
	homePath       string
	botLiteMode    bool
	logFilePath    string
}

type homeDirOption struct{ path string }

func HomeDir(path string) ClientOption               { return &homeDirOption{path: mustExpand(path)} }
func (o *homeDirOption) apply(c *clientConfig) error { c.homePath = o.path; return nil }

type executablePathOption struct{ path string }

func ExecutablePath(path string) ClientOption               { return &executablePathOption{path: mustExpand(path)} }
func (o *executablePathOption) apply(c *clientConfig) error { c.executablePath = o.path; return nil }

type botLiteOption struct{}

func BotLiteMode() ClientOption                      { return &botLiteOption{} }
func (o *botLiteOption) apply(c *clientConfig) error { c.botLiteMode = true; return nil }

type logFileOption struct{ path string }

func LogFilePath(path string) ClientOption           { return &logFileOption{path: mustExpand(path)} }
func (o *logFileOption) apply(c *clientConfig) error { c.logFilePath = o.path; return nil }
