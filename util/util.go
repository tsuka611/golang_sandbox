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

func ToFloat64(key string, val interface{}) (float64, error) {
	if v, ok := val.(float64); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(float64) -> %v", key, v))
	}
}

func ToString(key string, val interface{}) (string, error) {
	if v, ok := val.(string); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(string) -> %v", key, v))
	}
}
