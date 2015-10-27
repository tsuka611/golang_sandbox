package util

import (
	"encoding/json"
	"errors"
	"github.com/tsuka611/golang_sandbox/log"
	"reflect"
	"testing"
	"io/ioutil"
	"path/filepath"
)

func unmarchal(s string) map[string]interface{} {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	return m
}

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
	m := unmarchal(`{}`)
	if _, e := GetByFloat64(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByFloat64_illevalValueType_string(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": "TEST_VALUE"}`)
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
	m := unmarchal(`{"TEST_KEY": 11.1}`)
	expected := 11.1
	if v, e := GetByFloat64(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByFloat64_nevagiveValue(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": -11.1}`)
	expected := -11.1

	if v, e := GetByFloat64(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByString_keyNotExists(t *testing.T) {
	m := unmarchal(`{}`)
	if _, e := GetByString(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByString_illegalValueType_float(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": 1.1}`)
	if _, e := GetByString(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByString_empty(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": ""}`)
	expected := ""
	if v, e := GetByString(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByString_normal(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": "TEST_VALUE"}`)
	expected := "TEST_VALUE"
	if v, e := GetByString(m, "TEST_KEY"); e != nil {
		t.Errorf("error occurred. ERROR[%v]", e)
	} else if v != expected {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByStringArray_keyNotExists(t *testing.T) {
	m := unmarchal(`{}`)
	if _, e := GetByStringArray(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByStringArray_illegalValueType_string(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": "TEST_VALUE"}`)
	if _, e := GetByStringArray(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByStringArray_illegalValueType_in_int(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": ["x", "y", 1]}`)
	if _, e := GetByStringArray(m, "TEST_KEY"); e == nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	}
}

func TestGetByStringArray_empty(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": []}`)
	expected := []string{}
	if v, e := GetByStringArray(m, "TEST_KEY"); e != nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	} else if !reflect.DeepEqual(v, expected) {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestGetByStringArray_normal(t *testing.T) {
	m := unmarchal(`{"TEST_KEY": ["x", "y", "z"]}`)
	expected := []string{"x", "y", "z"}
	if v, e := GetByStringArray(m, "TEST_KEY"); e != nil {
		t.Errorf("expect error. VALUE[%v]", m["TEST_KEY"])
	} else if !reflect.DeepEqual(v, expected) {
		t.Errorf("expect `%v` but was `%v`.", expected, v)
	}
}

func TestExists_notExists(t *testing.T) {
	path := "xx"
	actual := Exists(path)
	expected := false
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestExists_fileExists(t *testing.T) {
	f, err := ioutil.TempFile("", "util_test")
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	path := f.Name()
	actual := Exists(path)
	expected := true
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestExists_dirExists(t *testing.T) {
	f, err := ioutil.TempFile("", "util_test")
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	path := filepath.Dir(f.Name())
	actual := Exists(path)
	expected := true
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestNewDir_fileExists(t *testing.T) {
	f, err := ioutil.TempFile("", "util_test")
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	expected := filepath.Dir(f.Name())
	if !Exists(expected) {
		t.Errorf("Cannot create Temp File.")
	}

	actual, err := NewDir(filepath.Dir(expected), filepath.Base(expected))
	if err == nil {
		t.Errorf(`Error must occur for [%v].`, err)
	}
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestNewDir_dirExists(t *testing.T) {
	f, err := ioutil.TempDir("", "util_test")
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	expected := f
	if !Exists(expected) {
		t.Errorf("Cannot create Temp Dir.")
	}

	actual, err := NewDir(filepath.Dir(expected), filepath.Base(expected))
	if err == nil {
		t.Errorf(`Error must occur for [%v].`, expected)
	}
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
}

func TestNewDir_success(t *testing.T) {
	f, err := ioutil.TempDir("", "util_test")
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	expected := filepath.Join(f, "success")
	if Exists(expected) {
		t.Errorf("Cannot create Temp Dir already exists.")
	}

	actual, err := NewDir(filepath.Dir(expected), filepath.Base(expected))
	if err != nil {
		t.Errorf("error occurred. ERROR[%v]", err)
	}
	if actual != expected {
		t.Errorf("expect `%v`, but was `%v`.", expected, actual)
	}
	if !Exists(actual) {
		t.Errorf("Cannot create New Dir.")
	}
}

func init() {
	log.SetLogLevel(log.L_TRACE)
}
