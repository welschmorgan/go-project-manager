package exec

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/welschmorgan/go-release-manager/config"
)

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
	stdoutPipe  io.ReadCloser
	stderrPipe  io.ReadCloser
	stdoutLines []string
	stderrLines []string
}

func NewCommand(name string, args ...string) *Command {
	return &Command{
		Cmd:         exec.Command(name, args...),
		name:        name,
		args:        args,
		stdoutLines: []string{},
		stderrLines: []string{},
	}
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
			log.Printf("Error executing command: %s......\n", err.Error())
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

func (c *Command) Run(filter CommandOutputFilter) error {
	if config.Get().Verbose || config.Get().DryRun {
		argStr := c.ArgString()
		fmt.Printf("%s* exec: %q %s\n", strings.Repeat("\t", config.Get().Indent), c.name, argStr)
	}
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
				log.Printf("Failed to wait for command, killing now... %v", c.Process.Kill())
				return err
			}
		}

		lines := map[string]bool{}

		defer c.stdoutPipe.Close()
		defer c.stderrPipe.Close()

		c.stdoutLines = []string{}
		c.stderrLines = []string{}

		for line, err := bufStdout.ReadString('\n'); err == nil; line, err = bufStdout.ReadString('\n') {
			if config.Get().Verbose {
				fmt.Fprintf(os.Stdout, "%s", line)
			}
			if filter&CommandOutputTrim != CommandOutputNoFilter {
				line = strings.TrimSpace(line)
			}
			if len(line) > 0 || filter&CommandOutputKeepEmpty == CommandOutputNoFilter {
				if filter&CommandOutputUnique != CommandOutputNoFilter {
					if ok := lines[line]; !ok {
						lines[line] = true
						c.stdoutLines = append(c.stdoutLines, line)
					}
				} else {
					c.stdoutLines = append(c.stdoutLines, line)
				}
			}
		}

		for line, err := bufStderr.ReadString('\n'); err == nil; line, err = bufStderr.ReadString('\n') {
			if config.Get().Verbose {
				fmt.Fprintf(os.Stderr, "%s", line)
			}
			if filter&CommandOutputTrim != CommandOutputNoFilter {
				line = strings.TrimSpace(line)
			}
			if filter&CommandOutputUnique != CommandOutputNoFilter {
				if ok := lines[line]; !ok {
					lines[line] = true
					c.stderrLines = append(c.stderrLines, line)
				}
			} else {
				c.stderrLines = append(c.stderrLines, line)
			}
		}
	}
	return nil
}

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func RunCommand(name string, args ...string) (exitCode int, stdout []string, stderr []string, err error) {
	cmd := NewCommand(name, args...)
	err = cmd.Run(CommandOutputTrim | CommandOutputUnique)
	return cmd.ExitCode(), cmd.StdoutLines(), cmd.StderrLines(), err
}

func DumpCommandErrors(exitCode int, errs []string) {
	level := ""
	color := ""
	indent := strings.Repeat("\t", config.Get().Indent)
	if exitCode != 0 {
		level = "error"
		color = "\033[1;31m"
	} else {
		level = "warning"
		color = "\033[1;33m"
	}
	shouldPrint := level == "error" || config.Get().Verbose
	if !shouldPrint {
		return
	}
	if SHOW_ERRORS && len(errs) > 0 {
		if len(errs) == 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s%s: %v\n", indent, color, level, "\033[0m", errs[0])
		} else {
			errStr := ""
			numErrs := 0
			for _, err := range errs {
				if len(strings.TrimSpace(err)) > 0 {
					if len(errStr) > 0 {
						errStr += "\n"
					}
					errStr += fmt.Sprintf("%s\t- %s", indent, strings.TrimSpace(err))
					numErrs += 1
				}
			}
			fmt.Fprintf(os.Stderr, "%s%s%d %s(s)%s:\n%s\n", indent, color, len(errs), level, "\033[0m", errStr)
		}
	}
}
