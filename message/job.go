package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/util"
	"strings"
	"github.com/tsuka611/golang_sandbox/job"
)

type Job struct {
	ID      job.JobID `json:id`      // JobID
	Command string    `jdon:command` // command search for $PATH
	Args    []string  `json:args`    // argument for command
	BaseEnv []string  `json:baseenv` // env values for common setting. (key=value pairs)
	Env     []string  `json:env`     // env values for custom setting. (key=value pairs)
	Dir     string    `json:dir`     // working directory for execute
}

func (e *Job) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("ID:%v", e.ID))
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

func (e *Job) UnmarshalJSON(b []byte) error {
	tmp := make(map[string]interface{})
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for key := range tmp {
		switch strings.ToLower(key) {
		case "id":
			if v, err := util.GetByString(tmp, key); err != nil {
				return err
			} else {
				e.ID = job.JobID(v)
			}
		case "command":
			if v, err := util.GetByString(tmp, key); err != nil {
				return err
			} else {
				e.Command = v
			}
		case "args":
			if v, err := util.GetByStringArray(tmp, key); err != nil {
				return err
			} else {
				e.Args = v
			}
		case "baseenv":
			if v, err := util.GetByStringArray(tmp, key); err != nil {
				return err
			} else {
				e.BaseEnv = v
			}
		case "env":
			if v, err := util.GetByStringArray(tmp, key); err != nil {
				return err
			} else {
				e.Env = v
			}
		case "dir":
			if v, err := util.GetByString(tmp, key); err != nil {
				return err
			} else {
				e.Dir = v
			}
		}
	}
	return nil
}

func NewJob(id job.JobID) *Job {
	return &Job{ID: id, Args: []string{}, BaseEnv: []string{}, Env: []string{}}
}
