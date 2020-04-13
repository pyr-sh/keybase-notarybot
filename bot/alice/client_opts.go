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

// Runs the service inside of the specified home directory, useful for running
// multiple services at once.
func HomeDir(path string) ClientOption               { return &homeDirOption{path: mustExpand(path)} }
func (o *homeDirOption) apply(c *clientConfig) error { c.homePath = o.path; return nil }

type executablePathOption struct{ path string }

// Uses the specified executable path to run all the Keybase commands.
func ExecutablePath(path string) ClientOption               { return &executablePathOption{path: mustExpand(path)} }
func (o *executablePathOption) apply(c *clientConfig) error { c.executablePath = o.path; return nil }

type botLiteOption struct{}

// Enables "bot lite mode", which turns off non-essential notifications such as typing.
func BotLiteMode() ClientOption                      { return &botLiteOption{} }
func (o *botLiteOption) apply(c *clientConfig) error { c.botLiteMode = true; return nil }

type logFileOption struct{ path string }

// If set, saves all the logs to a service-managed file if using the managed process mode
func LogFilePath(path string) ClientOption           { return &logFileOption{path: mustExpand(path)} }
func (o *logFileOption) apply(c *clientConfig) error { c.logFilePath = o.path; return nil }
