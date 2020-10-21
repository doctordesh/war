package war

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RunningState string

const (
	RunningStateNotStarted RunningState = "NOT_STARTED"
	RunningStateRunning    RunningState = "RUNNING"
	RunningStateStopped    RunningState = "STOPPED"
)

type Runnable interface {
	State() RunningState
	Start() error
	Stop() error
	ExitCode() (int, error)
}

type runnable struct {
	cmd *exec.Cmd

	state    RunningState
	exitCode int
}

func NewRunnable(cmd string, args []string, stdout, stderr io.Writer) Runnable {
	command := exec.Command(cmd, args...)

	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	command.Stdout = stdout
	command.Stderr = stderr

	return &runnable{state: RunningStateNotStarted, cmd: command}
}

func (self *runnable) State() RunningState {
	return self.state
}

func (self *runnable) Start() error {
	if self.state == RunningStateRunning {
		return fmt.Errorf("already started")
	}

	if self.state == RunningStateStopped {
		return fmt.Errorf("already done")
	}

	err := self.cmd.Start()
	if err != nil {
		return fmt.Errorf("could not start: %w", err)
	}

	go self.wait()

	self.state = RunningStateRunning

	return nil
}

func (self *runnable) Stop() error {
	if self.state == RunningStateRunning {
		fmt.Printf("Runnable kills process")
		err := self.cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("could not kill process: %w", err)
		}
	}

	self.state = RunningStateStopped

	return nil
}

func (self *runnable) ExitCode() (int, error) {
	if self.state != RunningStateStopped {
		return 0, fmt.Errorf("not done yet")
	}

	return self.exitCode, nil
}

func (self *runnable) wait() {
	err := self.cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Setting exit code: %+v\n", exiterr.ExitCode())
			self.exitCode = exiterr.ExitCode()
		} else {
			panic(err)
		}
	}

	self.state = RunningStateStopped
}
