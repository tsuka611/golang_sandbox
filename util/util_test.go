package util

import (
	"errors"
	"github.com/tsuka611/golang_sandbox/log"
	"testing"
)

func TestExtractOrPanic_panic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic not occurred.")
		}
	}()
	var actual interface{}
	var expected string = "OK"
	actual = ExtractOrPanic(func() (interface{}, error) {
		return expected, errors.New("Error")
	})

	if _, ok := actual.(string); !ok {
		t.Errorf("expect type assertion string [%v], but failed. ", actual)
	}
}

func TestExtractOrPanic_success(t *testing.T) {
	var actual interface{}
	var expected string = "OK"
	actual = ExtractOrPanic(func() (interface{}, error) {
		return expected, nil
	})
	if v, ok := actual.(string); !ok {
		t.Errorf("expect type assertion string [%v], but failed. ", actual)
	} else if v != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestExtractOrPanic_failTypeAssertion(t *testing.T) {
	var actual interface{}
	var expected int = 99
	actual = ExtractOrPanic(func() (interface{}, error) {
		return expected, nil
	})
	if _, ok := actual.(string); ok {
		t.Errorf("expect type assertion failed [%v], but success. ", actual)
	}
}

func TestGetByFloat64_keyNotExists(t *testing.T) {
	m := map[string]interface{}{}
	if _, e := GetByFloat64(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByFloat64_illevalValueType_string(t *testing.T) {
	m := map[string]interface{}{
		"TEST_KEY": "TEST_VALUE",
	}
	if _, e := GetByFloat64(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByFloat64_illegalValueType_integer(t *testing.T) {
	m := map[string]interface{}{
		"TEST_KEY": 1,
	}
	if _, e := GetByFloat64(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByFloat64_positiveValue(t *testing.T) {
	expected := 11.1
	m := map[string]interface{}{
		"TEST_KEY": expected,
	}
	if v, e := GetByFloat64(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByFloat64_nevagiveValue(t *testing.T) {
	expected := -11.1
	m := map[string]interface{}{
		"TEST_KEY": expected,
	}
	if v, e := GetByFloat64(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByString_keyNotExists(t *testing.T) {
	m := map[string]interface{}{}
	if _, e := GetByString(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByString_illegalValueType_float(t *testing.T) {
	m := map[string]interface{}{
		"TEST_KEY": 1.1,
	}
	if _, e := GetByString(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByString_empty(t *testing.T) {
	expected := ""
	m := map[string]interface{}{
		"TEST_KEY": expected,
	}
	if v, e := GetByString(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByString_normal(t *testing.T) {
	expected := "TEST_VALUE"
	m := map[string]interface{}{
		"TEST_KEY": expected,
	}
	if v, e := GetByString(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
