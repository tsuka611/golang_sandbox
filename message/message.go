package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/config"
	"github.com/tsuka611/golang_sandbox/log"
	"io"
)

type Operation string

const (
	OP_REGISTRATION = "REGISTRATION"
)

type Message struct {
	AppKey    config.AppKey
	Operation Operation
	Body      string
}

func NewMessage(appKey config.AppKey, op Operation, body interface{}) (*Message, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &Message{AppKey: appKey, Operation: op, Body: string(b)}, nil
}

func emptyMessage() *Message {
	return &Message{}
}

func SendMessage(w io.Writer, m *Message) (err error) {
	log.TRACE.Println("Try send message.")
	defer log.TRACE.Printlnf("Finish send message. Error[%v]", err)
	if _, err = io.WriteString(w, fmt.Sprintln(m.AppKey)); err != nil {
		return
	}
	if _, err = io.WriteString(w, fmt.Sprintln(m.Operation)); err != nil {
		return
	}
	_, err = io.WriteString(w, fmt.Sprintln(m.Body))
	return
}

func ReceiveMessage(r io.Reader) (m *Message, err error) {
	log.TRACE.Println("Try receive message.")
	defer log.TRACE.Printlnf("Finish receive message. Error[%v]", err)
	m = emptyMessage()
	sc := bufio.NewScanner(r)
	sc.Scan()
	m.AppKey = config.AppKey(sc.Text())
	sc.Scan()
	m.Operation = Operation(sc.Text())
	sc.Scan()
	m.Body = sc.Text()
	return m, sc.Err()
}
