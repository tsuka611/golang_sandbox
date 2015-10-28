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
	Port int
	Key  AppKey
}

func (c *agentConfig) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("coreConfig:%v", c.coreConfig))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("port:%v", c.Port))
	buf.WriteString(", ")
	switch {
	case len(c.Key) == 0:
		buf.WriteString(fmt.Sprintf("appKey:%v", ""))
	case len(c.Key) < 4:
		buf.WriteString(fmt.Sprintf("appKey:%v", "..."))
	default:
		buf.WriteString(fmt.Sprintf("appKey:%v", c.Key[0:3]+"..."))
	}
	return "{" + buf.String() + "}"
}

var (
	agentConf     *agentConfig
	onceAgentConf sync.Once
)

func AgentConfig() *agentConfig {
	onceAgentConf.Do(func() {
		agentConf = loadAgentConfig(CoreConfig())
	})
	return agentConf
}

func loadAgentConfig(core *coreConfig) *agentConfig {
	log.TRACE.Println("Start load AGENT config.")
	c := &agentConfig{coreConfig: core}

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
	if err = json.Unmarshal(file, c); err != nil {
		panic(err)
	}

	log.TRACE.Println("Fnish load AGENT config : ", c)
	return c
}
