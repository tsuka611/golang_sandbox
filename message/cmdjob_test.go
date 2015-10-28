package message

import (
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/job"
	"github.com/tsuka611/golang_sandbox/log"
	"reflect"
	"testing"
)

func TestString_cmdjob_empty(t *testing.T) {
	com := defCommon("first")
	base := newBaseJob(com, T_CMD, job.JobID("123"))
	actual := NewCmdJob(base, "ls").String()
	expected := fmt.Sprintf(`{Base:%v, Command:ls}`, base.String())
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func TestUnmarshalMarshal_cmdjob_normal(t *testing.T) {
	base := newBaseJob(defCommon("TestKey"), T_CMD, job.JobID("123"))
	base.AppKey = "xxxxx"
	base.Args = []string{"-l", "*"}
	base.BaseEnv = []string{"PATH=/bin"}
	base.Env = []string{"HOME=/home/hoge"}
	base.Dir = "/opt/local"
	expected := NewCmdJob(base, "TestCommand")

	actual := &CmdJob{}
	buf, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	log.INFO.Printlnf("Marshal JSON[%v]", string(buf))
	if err = json.Unmarshal(buf, actual); err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
