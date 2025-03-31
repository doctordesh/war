package war

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type RunnableTemplate struct {
	BinPath string
	Args    []string
	Env     []string
	Dir     string
	Stdout  io.Writer
	Stderr  io.Writer
}

func (self RunnableTemplate) Build() *runnable {
	cmd := &exec.Cmd{}
	cmd.Path = self.BinPath
	cmd.Args = self.Args
	cmd.Env = append(cmd.Environ(), self.Env...)
	cmd.Dir = self.Dir
	cmd.Stdout = self.Stdout
	cmd.Stderr = self.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return &runnable{state: RunningStateNotStarted, cmd: cmd}
}

type RunningState string

const (
	RunningStateNotStarted RunningState = "NOT_STARTED"
	RunningStateRunning    RunningState = "RUNNING"
	RunningStateStopped    RunningState = "STOPPED"
)

type runnable struct {
	cmd *exec.Cmd

	state    RunningState
	exitCode int
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
		err := syscall.Kill(-self.cmd.Process.Pid, syscall.SIGINT)
		if err != nil {
			err = syscall.Kill(-self.cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				return fmt.Errorf("could not kill process: %w", err)
			}
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
			// if exiterr.ExitCode() == -1 {
			// 	log.Printf("Runnable killed by war")
			// } else {
			// 	log.Printf("Setting exit code: %+v\n", exiterr.ExitCode())
			// }
			self.exitCode = exiterr.ExitCode()
		} else {
			panic(err)
		}
	}

	self.state = RunningStateStopped
}
