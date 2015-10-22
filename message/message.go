package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"github.com/tsuka611/golang_sandbox/util"
	"strings"
)

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

	for key := range tmp {
		switch strings.ToLower(key) {
		case "jobid":
			if v, err := util.GetByString(tmp, key); err != nil {
				return nil
			} else {
				m.JobID = JobID(v)
			}
		case "status":
			if v, err := util.GetByFloat64(tmp, key); err != nil {
				return nil
			} else {
				m.Status = Status(v)
			}
		case "exit":
			if v, err := util.GetByFloat64(tmp, key); err != nil {
				return nil
			} else {
				m.Exit = Exit(v)
			}
		}
	}

	return nil
}

func New(jobID JobID, status Status, exit Exit) *Message {
	return &Message{jobID, status, exit}
}

func Marshal(m *Message) (b []byte, e error) {
	b, e = json.Marshal(m)
	if e != nil {
		log.WARN.Printlnf("Marshal failed. Message[%v] Error[%v]", m, e)
	}
	return
}

func Unmarshal(b []byte) (m *Message, e error) {
	m = &Message{}
	e = json.Unmarshal(b, m)
	if e != nil {
		log.WARN.Printlnf("Unmarshal failed. String[%v] Error[%v]", string(b), e)
	}
	return
}
