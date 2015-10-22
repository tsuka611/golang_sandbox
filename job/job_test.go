package job

import (
	"bytes"
	"github.com/tsuka611/golang_sandbox/log"
	"os/exec"
	"testing"
	"time"
)

func TestRun_finishNormal(t *testing.T) {
	var out bytes.Buffer
	id := JobID("TestRun_finishNormal")
	cmd := exec.Command("sleep", "1")
	cmd.Stdout = &out
	cmd.Stderr = &out
	interrupt := make(chan bool)
	job := New(id, cmd, interrupt)
	finish := make(chan *Job)

	Run(job, finish)
	_ = <-finish
	log.TRACE.Println(out.String())
	if !job.Finished() {
		t.Error("expect Finished true, but was false.")
	}
	if !job.Success() {
		t.Error("expect Success true but was false.")
	}
	if job.ExitStatus() != 0 {
		t.Errorf("expect ExitStatus `%v` but was `%v`.", 0, job.ExitStatus())
	}
}

func TestRun_finishFailed(t *testing.T) {
	var out bytes.Buffer
	id := JobID("TestRun_finishFailed")
	cmd := exec.Command("cat", "NoSutchFile")
	cmd.Stdout = &out
	cmd.Stderr = &out
	interrupt := make(chan bool)
	job := New(id, cmd, interrupt)
	finish := make(chan *Job)

	Run(job, finish)
	_ = <-finish
	log.TRACE.Println(out.String())
	if !job.Finished() {
		t.Error("expect Finished true, but was false.")
	}
	if job.Success() {
		t.Error("expect Success false but was true.")
	}
	if job.ExitStatus() == 0 {
		t.Errorf("expect ExitStatus `%v` but was `%v`.", 0, job.ExitStatus())
	}
}

func TestRun_interrupt(t *testing.T) {
	var out bytes.Buffer
	id := JobID("TestRun_interrupt")
	cmd := exec.Command("sleep", "10")
	cmd.Stdout = &out
	cmd.Stderr = &out
	interrupt := make(chan bool)
	job := New(id, cmd, interrupt)
	finish := make(chan *Job)

	Run(job, finish)
	go func() {
		time.Sleep(2 * time.Second)
		close(interrupt)
	}()
	_ = <-finish
	log.TRACE.Println(out.String())
	if !job.Finished() {
		t.Error("expect Finished true, but was false.")
	}
	if job.Success() {
		t.Error("expect Success false but was true.")
	}
	if job.ExitStatus() == 0 {
		t.Errorf("expect ExitStatus `%v` but was `%v`.", 0, job.ExitStatus())
	}
}

func TestRun_commandNotFound(t *testing.T) {
	var out bytes.Buffer
	id := JobID("TestRun_commandNotFound")
	cmd := exec.Command("NoCmd")
	cmd.Stdout = &out
	cmd.Stderr = &out
	interrupt := make(chan bool)
	job := New(id, cmd, interrupt)
	finish := make(chan *Job)

	Run(job, finish)
	_ = <-finish
	log.TRACE.Println(out.String())
	if !job.Finished() {
		t.Error("expect Finished true, but was false.")
	}
	if job.Success() {
		t.Error("expect Success false but was true.")
	}
	if job.ExitStatus() == 0 {
		t.Errorf("expect ExitStatus `%v` but was `%v`.", 0, job.ExitStatus())
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
