package util

func ExtractOrPanic(f func() (interface{}, error)) interface{} {
	if val, err := f(); err != nil {
		panic(err)
	} else {
		return val
	}
}
