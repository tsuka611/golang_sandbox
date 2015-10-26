package message

import (
	"encoding/json"
	"github.com/tsuka611/golang_sandbox/config"
	"reflect"
	"testing"
)

func TestString_normalCase(t *testing.T) {
	actual := NewCommon(config.AppKey("xxxx")).String()
	expected := `{AppKey:xxxx}`
	if actual != expected {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}
}

func TestMarshalJSON_normalCase(t *testing.T) {
	buf, err := json.Marshal(NewCommon(config.AppKey(`"<>Hello`)))
	if err != nil {
		t.Error(err)
	}
	actual := string(buf)
	expected := `{"AppKey":"\"\u003c\u003eHello"}`
	if actual != expected {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}
}

func TestMarshalJSON_illegalType_integer(t *testing.T) {
	buf := []byte(`{"AppKey": 123}`)
	actual := NewCommon(config.AppKey("xxxx"))
	err := json.Unmarshal(buf, actual)
	if err == nil {
		t.Errorf(`Error must occur for [%v].`, string(buf))
	}
}

func TestUnmarshalJSON_normalCase(t *testing.T) {
	buf := []byte(`{"AppKey":"hello"}`)
	actual := NewCommon(config.AppKey("xxxx"))
	err := json.Unmarshal(buf, actual)
	if err != nil {
		t.Error(err)
	}
	expected := NewCommon(config.AppKey("hello"))
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}
}
