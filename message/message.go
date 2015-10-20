package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"strings"
)

type JobID string
type Status int
type Exit int
type Message struct {
	JobID  JobID  `json:jobid`
	Status Status `json:status`
	Exit   Exit   `json:exit`
}

func (m *Message) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("JobID:%v ", m.JobID))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Status:%v ", m.Status))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Exit:%v ", m.Exit))
	return "{" + buf.String() + "}"
}

func (m *Message) UnmarshalJSON(b []byte) error {
	tmp := make(map[string]interface{})
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for key, val := range tmp {
		switch strings.ToLower(key) {
		case "jobid":
			if v, err := convertString(key, val); err != nil {
				return nil
			} else {
				m.JobID = JobID(v)
			}
		case "status":
			if v, err := convertFloat64(key, val); err != nil {
				return nil
			} else {
				m.Status = Status(v)
			}
		case "exit":
			if v, err := convertFloat64(key, val); err != nil {
				return nil
			} else {
				m.Exit = Exit(v)
			}
		}
	}

	return nil
}

func convertFloat64(key string, val interface{}) (float64, error) {
	if v, ok := val.(float64); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(float64) -> %v", key, v))
	}
}

func convertString(key string, val interface{}) (string, error) {
	if v, ok := val.(string); ok {
		return v, nil
	} else {
		return v, errors.New(fmt.Sprintf("Cannot parse `%v`(string) -> %v", key, v))
	}
}

func New(jobID JobID, status Status, exit Exit) *Message {
	return &Message{jobID, status, exit}
}

func Marshal(m *Message) (b []byte, e error) {
	b, e = json.Marshal(m)
	if e != nil {
		log.WARN.Println("Marshal failed.", m, e)
	}
	return
}

func Unmarshal(b []byte) (m *Message, e error) {
	m = &Message{}
	e = json.Unmarshal(b, m)
	if e != nil {
		log.WARN.Println("Unmarshal failed.", string(b), e)
	}
	return
}
