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
	allowAutoFork  bool

	kbfsEnabled        bool
	kbfsExecutablePath string
	kbfsLogFilePath    string
	mountType          string
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

type allowAutoForkOption struct{ allow bool }

// If set, Start() starts a library-managed KBFS service
func AllowAutoFork(allow bool) ClientOption                { return &allowAutoForkOption{allow: allow} }
func (o *allowAutoForkOption) apply(c *clientConfig) error { c.allowAutoFork = o.allow; return nil }

type kbfsEnabledOption struct{ enabled bool }

// If set, Start() starts a library-managed KBFS service
func KBFSEnabled(enabled bool) ClientOption              { return &kbfsEnabledOption{enabled: enabled} }
func (o *kbfsEnabledOption) apply(c *clientConfig) error { c.kbfsEnabled = o.enabled; return nil }

type kbfsExecutablePathOption struct{ path string }

// Uses the specified executable path to run the KBFS service
func KBFSExecutablePath(path string) ClientOption { return &kbfsExecutablePathOption{path: path} }
func (o *kbfsExecutablePathOption) apply(c *clientConfig) error {
	c.kbfsExecutablePath = o.path
	return nil
}

type kbfsLogFileOption struct{ path string }

// If set, saves all the logs to a service-managed file if using the managed process mode
func KBFSLogFilePath(path string) ClientOption           { return &kbfsLogFileOption{path: mustExpand(path)} }
func (o *kbfsLogFileOption) apply(c *clientConfig) error { c.kbfsLogFilePath = o.path; return nil }

type mountTypeOption struct{ kind string }

// If set, the library-managed KBFS service receives the value as the mount-type flag (defaults to `none`)
func MountType(kind string) ClientOption               { return &mountTypeOption{kind: kind} }
func (o *mountTypeOption) apply(c *clientConfig) error { c.mountType = o.kind; return nil }
