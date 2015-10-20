package config

import (
	"bytes"
	"fmt"
	"github.com/tsuka611/golang_sandbox/log"
	"github.com/tsuka611/golang_sandbox/util"
	"os"
	"path/filepath"
	"sync"
	"flag"
)

type RunType int8
type coreConfig struct {
	prgDir string
	appDir string
	configFile string
}

const (
	_ = iota
	MASTER
	AGENT
	ADMIN
)

var (
	coreConf *coreConfig
	onceCoreConf sync.Once
	argConfFilePath string
	isTrace bool
)

func (c *coreConfig) PrgDir() string {
	return c.prgDir
}

func (c *coreConfig) AppDir() string {
	return c.appDir
}

func (c * coreConfig) ConfigFile() string {
	return c.configFile
}

func (c *coreConfig) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("prgDir:%v ", c.prgDir))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("appDir:%v ", c.appDir))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("configFile:%v ", c.configFile))
	return "{" + buf.String() + "}"
}

func CoreConfig() *coreConfig {
	onceCoreConf.Do(func() {
		loadCoreConfig()
	})
	return coreConf
}

func loadCoreConfig() {
	log.TRACE.Println("Start load core config.")
	c := &coreConfig{}
	c.prgDir = util.ExtractOrPanic(func() (interface{}, error) {
		return filepath.Abs(filepath.Dir(os.Args[0]))
	}).(string)
	c.appDir = util.ExtractOrPanic(func() (interface{}, error) {
		return filepath.Abs(filepath.Join(c.prgDir, ".."))
	}).(string)
	c.configFile = argConfFilePath
	coreConf = c
	log.TRACE.Println("Fnish load core config : ", coreConf)
}

func init() {
	flag.StringVar(&argConfFilePath, "c", "", "option for config file path.")
	flag.BoolVar(&isTrace, "trace", false, "option for output trace log.")
	flag.Parse()
	if (isTrace) {
		log.SetLogLevel(log.L_TRACE)
	}
}