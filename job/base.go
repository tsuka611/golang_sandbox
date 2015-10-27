package job

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"github.com/tsuka611/golang_sandbox/util"
	"os"
	"os/exec"
	"path/filepath"
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
	logFile    *os.File
	logWriter  *bufio.Writer
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

func (e *baseJob) Run() (err error) {
	if e.Finished() {
		err = errors.New(fmt.Sprintf("Job[%v] is already finifhed."))
		return
	}
	if e.alreadyStarted() {
		err = errors.New(fmt.Sprintf("Job[%v] is already started."))
		return
	}

	e.mu.Lock()
	e.started = true
	if e.logFile == nil {
		if e.logFile, err = buildLogfile(e.workDir); err != nil {
			return
		}
	}
	if e.logWriter == nil {
		e.logWriter = bufio.NewWriter(e.logFile)
	}
	if e.cmd.Stdout == nil {
		e.cmd.Stdout = e.logWriter
	}
	if e.cmd.Stderr == nil {
		e.cmd.Stderr = e.logWriter
	}

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

		if e.logWriter != nil {
			if err := e.logWriter.Flush(); err != nil {
				log.WARN.Printlnf("Log buffer flush failed. %v", err)
			}
		}
		if e.logFile != nil {
			if err := e.logFile.Close(); err != nil {
				log.WARN.Println("Log file close failed. %v", err)
			}
		}

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

func buildLogfile(dir string) (*os.File, error) {
	if !util.Exists(dir) {
		return nil, errors.New(fmt.Sprintf("%v is not a directory.", dir))
	}
	f, err := filepath.Abs(filepath.Join(dir, "output.log"))
	if err != nil {
		return nil, err
	}
	ret, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
