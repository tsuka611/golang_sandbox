package message

import (
	"bytes"
	"fmt"
	"github.com/tsuka611/golang_sandbox/config"
)

type Common struct {
	AppKey config.AppKey `json:appkey`
}

func (e *Common) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("AppKey:%v", e.AppKey))

	return "{" + buf.String() + "}"
}

func NewCommon(appKey config.AppKey) *Common {
	return &Common{AppKey: appKey}
}
