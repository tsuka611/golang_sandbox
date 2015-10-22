package message

import (
	"encoding/json"
	"github.com/tsuka611/golang_sandbox/log"
	"reflect"
	"testing"
)

func TestString_empty(t *testing.T) {
	actual := NewJob(JobID("123")).String()
	expected := `{ID:123, Command:, Args:[], BaseEnv:[], Env:[], Dir:}`
	if actual != expected {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func TestUnmarshalJSON_empty(t *testing.T) {
	org := NewJob(JobID("123"))
	actual := NewJob(JobID("x"))

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

func TestUnmarshalJSON_normal(t *testing.T) {
	org := NewJob(JobID("123"))
	org.Command = "ls"
	org.Args = []string{"-l", "*"}
	org.BaseEnv = []string{"PATH=/bin"}
	org.Env = []string{"HOME=/home/hoge"}
	org.Dir = "/opt/local"

	actual := NewJob(JobID("x"))

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

func init() {
	log.SetLogLevel(log.L_TRACE)
}
