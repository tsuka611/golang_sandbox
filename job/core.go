package job

import (
	"errors"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"os/exec"
	"sync"
	"syscall"
)

type JobID string
type ExitStatus int

type Runner interface {
	ExitStatus() (ExitStatus, error)
	Success() bool
	Finished() bool
	Run() error
	Wait() error
}

type baseJob struct {
	JobID      JobID
	cmd        *exec.Cmd
	exit       chan error
	interrupt  <-chan bool
	exitStatus syscall.WaitStatus
	finished   bool
	started    bool
	workDir    string
	mu         sync.Mutex
	once       sync.Once
}

func (e *baseJob) ExitStatus() (ExitStatus, error) {
	if !e.Finished() {
		return ExitStatus(-1), errors.New("Job is not finished.")
	}
	return ExitStatus(e.exitStatus.ExitStatus()), nil
}

func (e *baseJob) Success() bool {
	if !e.Finished() {
		return false
	}
	sts, _ := e.ExitStatus()
	return sts == 0
}

func (e *baseJob) Finished() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.finished
}

func (e *baseJob) Run() error {
	if e.Finished() {
		return errors.New(fmt.Sprintf("Job[%v] is already finifhed."))
	}
	if e.alreadyStarted() {
		return errors.New(fmt.Sprintf("Job[%v] is already started."))
	}
	e.mu.Lock()
	e.started = true
	e.mu.Unlock()

	log.TRACE.Printlnf("Running job. %v", e.JobID)
	if err := e.cmd.Start(); err != nil {
		e.closeJob()
		return err
	}

	go func() {
		defer close(e.exit)
		e.exit <- e.cmd.Wait()
	}()

	log.TRACE.Printlnf("Start Running job success. %v", e.JobID)
	return nil
}

func (e *baseJob) Wait() (err error) {
	if e.Finished() {
		err = errors.New("Job is already finished.")
		return
	}
	if !e.alreadyStarted() {
		err = errors.New("Job is not started.")
		return
	}

	defer e.closeJob()

	for {
		select {
		case err, _ = <-e.exit:
			if err != nil {
				log.ERROR.Printlnf("Job Error exists. Err[%v] Job[%v]", err, e)
			}
			return
		case _, _ = <-e.interrupt:
			log.WARN.Printlnf("Interrupt message received. Job[%v]", e)
			if err := e.cmd.Process.Kill(); err != nil {
				log.ERROR.Printlnf("Error while Job kill. Error[%v] Job[%v]", err, e)
			}
			_, _ = <-e.exit
			return
		}
	}
	return
}

func (e *baseJob) closeJob() {
	e.once.Do(func() {
		log.TRACE.Printlnf("Start closing job. %v", e)
		e.mu.Lock()
		defer e.mu.Unlock()

		e.finished = true

		ps := e.cmd.ProcessState
		if ps == nil {
			log.TRACE.Printlnf("Cannot get ProcessStatus, may be fail running. [%v]", e)
			return
		}
		sys := ps.Sys()
		if sys == nil {
			panic(errors.New(fmt.Sprintf("Cannot get ProcessStatus sys. [%v]", e)))
		}
		if sts, ok := sys.(syscall.WaitStatus); ok {
			e.exitStatus = sts
			return
		}
		panic(fmt.Sprintf("Cannot extract WaitStatus from [%v].", sys))
	})
}

func (e *baseJob) alreadyStarted() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.started
}

func newBaseJob(jobID JobID, cmd *exec.Cmd, interrupt <-chan bool, workDir string) *baseJob {
	return &baseJob{
		JobID:      jobID,
		cmd:        cmd,
		exit:       make(chan error),
		interrupt:  interrupt,
		exitStatus: syscall.WaitStatus(0x7F00),
		finished:   false,
		workDir:    workDir,
	}
}
