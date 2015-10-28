package message

import (
	"bytes"
	"fmt"
)

type CmdJob struct {
	*baseJob
	Command string `jdon:command` // command search for $PATH
}

func NewCmdJob(b *baseJob, c string) *CmdJob {
	return &CmdJob{baseJob: b, Command: c}
}

func (e *CmdJob) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("Base:%v", e.baseJob))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Command:%v", e.Command))

	return "{" + buf.String() + "}"
}
