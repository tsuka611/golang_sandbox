package job

import (
	"github.com/tsuka611/golang_sandbox/log"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func testCmdJob(name string, arg ...string) (*baseJob, chan bool) {
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

func TestCmdJobFinished_notRunJob(t *testing.T) {
	job, _ := testCmdJob("sleep", "1")
	actual := job.Finished()
	expected := false
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestCmdJobFinished_finishedJob(t *testing.T) {
	job, _ := testCmdJob("sleep", "3")

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

func TestCmdJobFinished_interruptedJob(t *testing.T) {
	job, interrupt := testCmdJob("sleep", "10")

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

func TestCmdJobExitStatus_notRunJob(t *testing.T) {
	job, _ := testCmdJob("sleep", "10")
	actual, err := job.ExitStatus()
	expected := ExitStatus(-1)
	if err == nil {
		t.Errorf(`Error must occur for [%v].`, job)
	}
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestCmdJobExitStatus_finishedJob(t *testing.T) {
	job, _ := testCmdJob("sleep", "1")
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

func TestCmdJobExitStatus_interruptedJob(t *testing.T) {
	job, interrupt := testCmdJob("sleep", "10")
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

func TestCmdJobExitStatus_notCmdJob(t *testing.T) {
	job, _ := testCmdJob("xxsleep", "10")
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

func TestCmdJobExitStatus_failCmdJob(t *testing.T) {
	job, _ := testCmdJob("ping")
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

func TestCmdJobRun_outputLogCheck_stdout(t *testing.T) {
	job, _ := testCmdJob("echo", "SampleOutputData")
	job.cmd.Stdout = nil

	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if err := job.Wait(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if sts, err := job.ExitStatus(); err != nil || sts != 0 {
		t.Errorf("Command exec failed. STATUS[%v] ERROR[%v]", sts, err)
	}

	expected := "SampleOutputData"
	buf, err := ioutil.ReadFile(job.logFile.Name())
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	actual := strings.TrimSpace(string(buf))
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, actual)
	}
}

func TestCmdJobRun_outputLogCheck_stderr(t *testing.T) {
	job, _ := testCmdJob("ping", "xx")
	job.cmd.Stderr = nil

	if err := job.Run(); err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if err := job.Wait(); err == nil {
		t.Errorf(`Error must occur for [%v].`, job)
	}
	if sts, err := job.ExitStatus(); err != nil || sts == 0 {
		t.Errorf("Command exec success. STATUS[%v] ERROR[%v]", sts, err)
	}

	buf, err := ioutil.ReadFile(job.logFile.Name())
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	actual := strings.TrimSpace(string(buf))
	if len(actual) < 1 {
		t.Errorf("expect `%v` but was `%v`.", "NotEmptyString", actual)
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
