package message

import (
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/config"
	"github.com/tsuka611/golang_sandbox/job"
	"github.com/tsuka611/golang_sandbox/log"
	"reflect"
	"testing"
)

func defCommon(s string) *Common {
	return NewCommon(config.AppKey(s))
}

func TestString_basejob_empty(t *testing.T) {
	com := defCommon("first")
	actual := newBaseJob(com, T_CMD, job.JobID("123")).String()
	expected := fmt.Sprintf(`{Common:%v, ID:123, Args:[], BaseEnv:[], Env:[], Dir:}`, com.String())
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func TestUnmarshalMarshal_basejob_empty(t *testing.T) {
	org := newBaseJob(defCommon("TestKey"), T_CMD, job.JobID("123"))
	actual := newBaseJob(defCommon("x"), T_CMD, job.JobID("x"))

	buf, err := json.Marshal(org)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	log.INFO.Printlnf("Marshal JSON[%v]", string(buf))
	if err = json.Unmarshal(buf, actual); err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if !reflect.DeepEqual(org, actual) {
		t.Errorf("expect `%v` but was `%v`", org, actual)
	}
}

func TestUnmarshalMarshal_basejob_normal(t *testing.T) {
	org := newBaseJob(defCommon("TestKey"), T_CMD, job.JobID("123"))
	org.AppKey = "xxxxx"
	org.Args = []string{"-l", "*"}
	org.BaseEnv = []string{"PATH=/bin"}
	org.Env = []string{"HOME=/home/hoge"}
	org.Dir = "/opt/local"

	actual := newBaseJob(defCommon("x"), T_CMD, job.JobID("x"))

	buf, err := json.Marshal(org)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	log.INFO.Printlnf("Marshal JSON[%v]", string(buf))
	if err = json.Unmarshal(buf, actual); err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if !reflect.DeepEqual(org, actual) {
		t.Errorf("expect `%v` but was `%v`", org, actual)
	}
}

func TestUnmarshalJSON_basejob_empty(t *testing.T) {
	src := `{"AppKey":"xxxx", "JobID":"123"}`
	actual := newBaseJob(defCommon("x"), T_CMD, job.JobID("x"))

	err := json.Unmarshal([]byte(src), actual)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	expected := newBaseJob(defCommon("xxxx"), T_CMD, job.JobID("123"))
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func TestUnmarshalJSON_basejob_normal(t *testing.T) {
	src := `{"AppKey":"xxxx", "JobID":"123", "Command":"ls", "Args":["*"], "BaseEnv":["BE1=111"], "Env":["E2=22"], "Dir":"TestDir"}`
	actual := newBaseJob(defCommon("x"), T_CMD, job.JobID("x"))

	err := json.Unmarshal([]byte(src), actual)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	expected := newBaseJob(defCommon("x"), T_CMD, job.JobID("x"))
	expected.AppKey = "xxxx"
	expected.JobID = "123"
	expected.Args = []string{"*"}
	expected.BaseEnv = []string{"BE1=111"}
	expected.Env = []string{"E2=22"}
	expected.Dir = "TestDir"

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
