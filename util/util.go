package util

import (
	"errors"
	"fmt"
)

func ExtractOrPanic(f func() (interface{}, error)) interface{} {
	if val, err := f(); err != nil {
		panic(err)
	} else {
		return val
	}
}

func GetByFloat64(m map[string]interface{}, key string) (float64, error) {
	val := m[key]
	if v, ok := val.(float64); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(float64) -> %v", key, val))
	}
}

func GetByString(m map[string]interface{}, key string) (string, error) {
	val := m[key]
	if v, ok := val.(string); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(string) -> %v", key, val))
	}
}

func GetByStringArray(m map[string]interface{}, key string) ([]string, error) {
	val := m[key]
	var a []interface{}
	var ok bool
	if a, ok = val.([]interface{}); !ok {
		return []string{}, errors.New(fmt.Sprintf("Cannot parse `%v`([]interface) -> %v", key, val))
	}
	ret := make([]string, len(a))
	for i, e := range a {
		var v string
		if v, ok = e.(string); !ok {
			return []string{}, errors.New(fmt.Sprintf("Cannot parse in array `%v` at index:%v", e, i))
		}
		ret[i] = v
	}
	return ret, nil
}
