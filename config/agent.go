package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"github.com/tsuka611/golang_sandbox/util"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type agentConfig struct {
	*coreConfig
	port int
}

func (c *agentConfig) Port() int {
	return c.port
}

func (c *agentConfig) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("coreConfig:%v ", c.coreConfig))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("port:%v ", c.port))
	return "{" + buf.String() + "}"
}

var (
	agentConf     *agentConfig
	onceAgentConf sync.Once
)

func AgentConfig() *agentConfig {
	onceAgentConf.Do(func() {
		loadAgentConfig()
	})
	return agentConf
}

func loadAgentConfig() {
	log.TRACE.Println("Start load AGENT config.")
	c := &agentConfig{CoreConfig(), 0}

	confFile := c.ConfigFile()
	if len(confFile) == 0 {
		confFile = util.ExtractOrPanic(func() (interface{}, error) {
			return filepath.Abs(filepath.Join(c.AppDir(), "conf", "agent_conf.json"))
		}).(string)
	}
	log.TRACE.Println("Try to load config file: ", confFile)
	file, err := ioutil.ReadFile(confFile)
	if err != nil {
		panic(err)
	}

	log.TRACE.Println("Try to parse config file.")
	m := make(map[string]interface{})
	if err = json.Unmarshal(file, &m); err != nil {
		panic(err)
	}

	if val, ok := m["port"].(float64); !ok {
		panic(fmt.Sprintf("Cannot get `port` value. [%v]", m["port"]))
	} else {
		c.port = int(val)
	}

	agentConf = c
	log.TRACE.Println("Fnish load AGENT config : ", agentConf)
}
