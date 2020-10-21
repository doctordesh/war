package war

import "fmt"

type Runner interface {
	Run(Runnable) error
}

type runner struct {
	lastRunnable Runnable
}

func NewRunner() *runner {
	return &runner{}
}

func (self *runner) Run(runnable Runnable) error {
	var err error

	if self.lastRunnable != nil {
		fmt.Printf("Runner triest to stop 'lastRunnable'")
		err = self.lastRunnable.Stop()
		_ = err
	}

	fmt.Printf("Runner starts new runnable")
	err = runnable.Start()
	if err != nil {
		return err
	}

	self.lastRunnable = runnable

	return nil
}
