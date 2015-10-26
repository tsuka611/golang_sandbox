package message

import (
	"bytes"
	"fmt"
	"github.com/tsuka611/golang_sandbox/job"
)

type Job struct {
	*Common
	JobID   job.JobID `json:id`      // JobID
	Command string    `jdon:command` // command search for $PATH
	Args    []string  `json:args`    // argument for command
	BaseEnv []string  `json:baseenv` // env values for common setting. (key=value pairs)
	Env     []string  `json:env`     // env values for custom setting. (key=value pairs)
	Dir     string    `json:dir`     // working directory for execute
}

func (e *Job) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("Common:%v", e.Common))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("ID:%v", e.JobID))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Command:%v", e.Command))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Args:%v", e.Args))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("BaseEnv:%v", e.BaseEnv))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Env:%v", e.Env))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Dir:%v", e.Dir))

	return "{" + buf.String() + "}"
}

func NewJob(c *Common, jobId job.JobID) *Job {
	return &Job{Common: c, JobID: jobId, Args: []string{}, BaseEnv: []string{}, Env: []string{}}
}
