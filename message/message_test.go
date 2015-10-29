package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSendReceive_binary_randByte(t *testing.T) {
	m, err := NewMessage("TestAppKey", OP_REGISTRATION, "TestBody")
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	var buf bytes.Buffer
	err = SendMessage(&buf, m)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}
	a, err := ReceiveMessage(&buf)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if m.AppKey != a.AppKey {
		t.Errorf("expect `%v` but was `%v`", m.AppKey, a.AppKey)
	}

	if m.Operation != a.Operation {
		t.Errorf("expect `%v` but was `%v`", m.Operation, a.Operation)
	}
	if m.Body != a.Body {
		t.Errorf("expect `%v` but was `%v`", m.Body, a.Body)
	}
}

func TestSendMessage_emptyMess(t *testing.T) {
	m := emptyMessage()
	var actual, expected bytes.Buffer

	expected.WriteString(fmt.Sprintln(""))
	expected.WriteString(fmt.Sprintln(""))
	expected.WriteString(fmt.Sprintln(""))

	if err := SendMessage(&actual, m); err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if bytes.Compare(actual.Bytes(), expected.Bytes()) != 0 {
		t.Errorf("expect `%v` but was `%v`", expected.String(), actual.String())
	}
}

func TestSendMessage_normal(t *testing.T) {
	m, err := NewMessage("TestAppKey", OP_REGISTRATION, "TestBody")
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}
	var actual, expected bytes.Buffer

	expected.WriteString(fmt.Sprintln("TestAppKey"))
	expected.WriteString(fmt.Sprintln(OP_REGISTRATION))
	if body, err := json.Marshal("TestBody"); err != nil {
		t.Errorf("Error occurred. %v", err)
	} else {
		expected.WriteString(fmt.Sprintln(string(body)))
	}

	if err := SendMessage(&actual, m); err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if bytes.Compare(actual.Bytes(), expected.Bytes()) != 0 {
		t.Errorf("expect `%v` but was `%v`", expected.String(), actual.String())
	}
}

func TestReceiveMessage_empty(t *testing.T) {
	var buf bytes.Buffer
	expected := emptyMessage()
	actual, err := ReceiveMessage(&buf)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}

func TestReceiveMessage_normal(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintln("TestAppKey"))
	buf.WriteString(fmt.Sprintln(OP_REGISTRATION))
	if body, err := json.Marshal("TestBody"); err != nil {
		t.Errorf("Error occurred. %v", err)
	} else {
		buf.WriteString(fmt.Sprintln(string(body)))
	}

	expected, err := NewMessage("TestAppKey", OP_REGISTRATION, "TestBody")
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	actual, err := ReceiveMessage(&buf)
	if err != nil {
		t.Errorf("Error occurred. %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expect `%v` but was `%v`", expected, actual)
	}
}
