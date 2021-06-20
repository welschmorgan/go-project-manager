package exec

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
)

type RunOptions struct {
	FailCommandWhenStderrContainsErrors bool
	TreatWarningsAsErrors               bool
}

var defaultRunOptions = RunOptions{
	FailCommandWhenStderrContainsErrors: true,
	TreatWarningsAsErrors:               false,
}

func DefaultRunOptions() RunOptions {
	return defaultRunOptions
}

var (
	SHOW_ERRORS = true
)

type CommandOutputFilter uint8

const (
	CommandOutputNoFilter  CommandOutputFilter = iota
	CommandOutputUnique    CommandOutputFilter = 1 << iota
	CommandOutputTrim      CommandOutputFilter = 1 << iota
	CommandOutputKeepEmpty CommandOutputFilter = 1 << iota
)

type Command struct {
	*exec.Cmd
	name        string
	args        []string
	options     RunOptions
	stdoutPipe  io.ReadCloser
	stderrPipe  io.ReadCloser
	stdoutLines []string
	stderrLines []string
}

func NewCommand(name string, args ...string) *Command {
	return NewCommandWithOptions(name, args, defaultRunOptions)
}

func NewCommandWithOptions(name string, args []string, options RunOptions) *Command {
	return &Command{
		Cmd:         exec.Command(name, args...),
		name:        name,
		options:     defaultRunOptions,
		args:        args,
		stdoutLines: []string{},
		stderrLines: []string{},
	}
}

func (c *Command) SetOptions(o RunOptions) {
	c.options = o
}

func (c *Command) Options() RunOptions {
	return c.options
}

func (c *Command) ArgString() string {
	argStr := ""
	for _, a := range c.args {
		if len(argStr) > 0 {
			argStr += " "
		}
		argStr += fmt.Sprintf("%q", a)
	}
	return argStr
}

func (c *Command) Start() error {
	if !config.Get().DryRun {
		c.stdoutPipe, _ = c.Cmd.StdoutPipe()
		c.stderrPipe, _ = c.Cmd.StderrPipe()
		if err := c.Cmd.Start(); err != nil {
			log.Printf("Cannot execute command: %s\n", err.Error())
			return err
		}
	}
	return nil
}

func (c *Command) ExitCode() int {
	return c.ProcessState.ExitCode()
}

func (c *Command) Stdout() string {
	return strings.Join(c.stdoutLines, "\n")
}

func (c *Command) Stderr() string {
	return strings.Join(c.stderrLines, "\n")
}

func (c *Command) StdoutLines() []string {
	return c.stdoutLines
}

func (c *Command) StderrLines() []string {
	return c.stderrLines
}

func (c *Command) readAll(r *bytes.Buffer, f CommandOutputFilter) []string {
	lines := map[string]bool{}
	retLines := []string{}
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		log.Traceln(line)
		if f&CommandOutputTrim != CommandOutputNoFilter {
			line = strings.TrimSpace(line)
		}
		if len(line) > 0 || f&CommandOutputKeepEmpty == CommandOutputNoFilter {
			if f&CommandOutputUnique != CommandOutputNoFilter {
				if ok := lines[line]; !ok {
					lines[line] = true
					retLines = append(retLines, line)
				}
			} else {
				retLines = append(retLines, line)
			}
		}
	}
	return retLines
}

func (c *Command) Run(filter CommandOutputFilter) error {
	argStr := c.ArgString()
	log.Tracef("%s* exec: %q %s\n", strings.Repeat("\t", config.Get().Indent), c.name, argStr)

	var wg sync.WaitGroup
	var bufStdout *bytes.Buffer = bytes.NewBufferString("")
	var bufStderr *bytes.Buffer = bytes.NewBufferString("")
	if !config.Get().DryRun {
		if err := c.Start(); err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			io.Copy(bufStdout, c.stdoutPipe)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			io.Copy(bufStderr, c.stderrPipe)
		}()

		wg.Wait()

		if err := c.Wait(); err != nil {
			if !c.ProcessState.Exited() {
				log.Errorf("Failed to wait for command, killing now... %v", c.Process.Kill())
				return err
			}
		}

		defer c.stdoutPipe.Close()
		defer c.stderrPipe.Close()

		c.stdoutLines = c.readAll(bufStdout, filter)
		c.stderrLines = c.readAll(bufStderr, filter|^CommandOutputUnique|^CommandOutputKeepEmpty)

		if len(c.stderrLines) > 0 && c.options.FailCommandWhenStderrContainsErrors && (c.ExitCode() != 0 || c.options.TreatWarningsAsErrors) {
			return fmt.Errorf("failed to run command '%s', %s", c.name, c.Stderr())
		}
	}

	return nil
}

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func RunCommand(name string, args ...string) (exitCode int, stdout []string, stderr []string, err error) {
	return RunCommandWithOptions(name, args, defaultRunOptions)
}

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func RunCommandWithOptions(name string, args []string, options RunOptions) (exitCode int, stdout []string, stderr []string, err error) {
	cmd := NewCommandWithOptions(name, args, options)
	err = cmd.Run(CommandOutputTrim | CommandOutputUnique)
	return cmd.ExitCode(), cmd.StdoutLines(), cmd.StderrLines(), err
}

func DumpCommandErrors(exitCode int, errs ...string) {
	DumpCommandErrorsWithOptions(exitCode, errs, defaultRunOptions)
}

func DumpCommandErrorsWithOptions(exitCode int, errs []string, options RunOptions) {
	level := ""
	color := ""
	if exitCode != 0 {
		level = "error"
		color = "\033[1;31m"
	} else {
		level = "warning"
		color = "\033[1;33m"
	}
	logErr := func(err string) {
		if level == "warning" && !options.TreatWarningsAsErrors {
			log.Warnf("%s%s: %v\n", color, "\033[0m", err)
		} else {
			log.Errorf("%s%s: %v\n", color, "\033[0m", err)
		}
	}
	if SHOW_ERRORS && len(errs) > 0 {
		if len(errs) == 1 {
			logErr(errs[0])
		} else {
			formattedErrs := []string{}
			numErrs := 0
			for _, err := range errs {
				if len(strings.TrimSpace(err)) > 0 {
					formattedErrs = append(formattedErrs, fmt.Sprintf("\t- %s", strings.TrimSpace(err)))
					numErrs += 1
				}
			}
			logErr(fmt.Sprintf("%d %s(s):\n%s\n", len(formattedErrs), level, strings.Join(formattedErrs, "\n")))
		}
	}
}
