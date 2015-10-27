package job

import (
	"github.com/tsuka611/golang_sandbox/log"
	"io/ioutil"
	"os/exec"
	"testing"
)

func testJob(name string,arg ...string) (*baseJob, chan bool) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = log.INFO
	cmd.Stderr = log.ERROR

	workDir, err := ioutil.TempDir("", "core_test")
	if err != nil {
		panic(err)
	}

	ch := make(chan bool)
	return newBaseJob(JobID("123x"), cmd, ch, workDir), ch
}

func TestFinished_notRunJob(t *testing.T) {
	job, _ := testJob("sleep", "1")
	actual := job.Finished()
	expected := false
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestFinished_finishedJob(t *testing.T) {
	job, _ := testJob("sleep", "3")

	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if err := job.Wait(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}

	actual := job.Finished()
	expected := true
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestFinished_interruptedJob(t *testing.T) {
	job, interrupt := testJob("sleep", "10")

	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	close(interrupt)
	if err := job.Wait(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}

	actual := job.Finished()
	expected := true
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestExitStatus_notRunJob(t *testing.T) {
	job,_ := testJob("sleep", "10")
	actual, err := job.ExitStatus()
	expected := ExitStatus(-1)
	if err == nil {
		t.Errorf(`Error must occur for [%v].`, job)
	}
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestExitStatus_finishedJob(t *testing.T) {
	job,_ := testJob("sleep", "1")
	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if err := job.Wait(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	actual, err := job.ExitStatus()
	expected := ExitStatus(0)
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestExitStatus_interruptedJob(t *testing.T) {
	job,interrupt := testJob("sleep", "10")
	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	close(interrupt)
	job.Wait()
	actual, err := job.ExitStatus()
	expected := ExitStatus(-1)
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestExitStatus_notCmdJob(t *testing.T) {
	job,_ := testJob("xxsleep", "10")
	if err := job.Run(); err == nil {
		t.Errorf(`Error must occur for [%v].`, job)
	}
	job.Wait()
	actual, err := job.ExitStatus()
	expected := ExitStatus(127)
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestExitStatus_failCmdJob(t *testing.T) {
	job,_ := testJob("ping")
	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if err := job.Wait(); err == nil {
		t.Errorf(`Error must occur for [%v].`, job)
	}
	actual, err := job.ExitStatus()
	expected := ExitStatus(0)
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if actual == expected {
		t.Errorf("expect not `%v` but was `%v`.", expected, actual)
	}
}



func init() {
	log.SetLogLevel(log.L_TRACE)
}
