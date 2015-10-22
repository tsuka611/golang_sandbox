package job

import (
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"os/exec"
	"sync"
	"syscall"
)

type Job struct {
	Id         string
	cmd        *exec.Cmd
	interrupt  <-chan bool
	isFinished bool
	waitStatus syscall.WaitStatus
	mu         sync.Mutex
}

func (j *Job) ExitStatus() int {
	return j.status().ExitStatus()
}

func (j *Job) status() syscall.WaitStatus {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.waitStatus
}

func (j *Job) setWaitStatus(w syscall.WaitStatus) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.waitStatus = w
}

func (j *Job) setFinished(finished bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.isFinished = finished
}

func (j *Job) updateStatus() {
	j.mu.Lock()
	defer j.mu.Unlock()

	ps := j.cmd.ProcessState
	if ps == nil {
		panic(fmt.Sprintf("Cannot get ProcessStatus. [%v]", j))
	}
	sys := ps.Sys()
	if sys == nil {
		panic(fmt.Sprintf("Cannot get ProcessStatus sys. [%v]", j))
	}
	if sts, ok := sys.(syscall.WaitStatus); ok {
		j.waitStatus = sts
		return
	}
	panic(fmt.Sprintf("Cannot extract WaitStatus from [%v].", sys))
}

func (j *Job) Finished() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.isFinished
}

func (j *Job) Success() bool {
	return j.Finished() && j.ExitStatus() == 0
}

func (j *Job) runCmd() <-chan bool {
	exit := make(chan bool)
	go func(ch chan bool) {
		defer close(ch)
		if err := j.cmd.Run(); err != nil {
			log.ERROR.Println(err)
		} else {
			j.updateStatus()
		}
		j.setFinished(true)
	}(exit)
	return exit
}

func New(Id string, cmd *exec.Cmd, interrupt <-chan bool) *Job {
	return &Job{Id: Id, cmd: cmd, interrupt: interrupt, waitStatus: syscall.WaitStatus(0x7F00)}
}

func Run(job *Job, finish chan<- *Job) {
	go runBackGround(job, finish)
}

func runBackGround(job *Job, finish chan<- *Job) {
	log.TRACE.Printlnf("Running job. Job[%v]", job)
	defer func() { finish <- job }()
	exit := job.runCmd()
	for {
		select {
		case _, _ = <-exit:
			log.TRACE.Printlnf("Exit job. Job[%v]", job)
			return
		case _, _ = <-job.interrupt:
			log.WARN.Printlnf("Interrupt message received. Job[%v]", job)
			if err := job.cmd.Process.Kill(); err != nil {
				log.ERROR.Printlnf("Error while Job kill. Error[%v] Job[%v]", err, job)
			}
			_, _ = <-exit
			return
		}
	}
}
