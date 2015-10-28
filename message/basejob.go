package message

import (
	"bytes"
	"fmt"
	"github.com/tsuka611/golang_sandbox/job"
)

type JobType int

const (
	T_CMD = JobType(1) << iota
)

type baseJob struct {
	*Common
	JobType JobType   `json:jobtype` // jobtype
	JobID   job.JobID `json:jobid`   // JobID
	Args    []string  `json:args`    // argument for command
	BaseEnv []string  `json:baseenv` // env values for common setting. (key=value pairs)
	Env     []string  `json:env`     // env values for custom setting. (key=value pairs)
	Dir     string    `json:dir`     // working directory for execute
}

func (e *baseJob) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("Common:%v", e.Common))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("ID:%v", e.JobID))
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

func newBaseJob(c *Common, jobType JobType, jobId job.JobID) *baseJob {
	return &baseJob{Common: c, JobType: jobType, JobID: jobId, Args: []string{}, BaseEnv: []string{}, Env: []string{}}
}
