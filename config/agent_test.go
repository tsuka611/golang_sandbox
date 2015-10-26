package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func tmpJson(s string) string {
	f, err := ioutil.TempFile("", "config_test")
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(s); err != nil {
		panic(err)
	}
	var ret string
	if ret, err = filepath.Abs(f.Name()); err != nil {
		panic(err)
	}
	return ret
}

func defCoreConf(prgDir, appDir, configFile string) *coreConfig {
	return &coreConfig{prgDir, appDir, configFile}
}

func TestString_appKeyEmpty(t *testing.T) {
	j := `
{
	"port": 123,
	"key": ""
}`
	c := loadAgentConfig(defCoreConf("prg/dir", "app/dir", tmpJson(j)))
	actual := c.String()
	expected := fmt.Sprintf(`{coreConfig:%v, port:123, appKey:}`, c.coreConfig.String())
	if actual != expected {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}

}

func TestString_appKey3letter(t *testing.T) {
	j := `
{
	"port": 123,
	"key": "987"
}`
	c := loadAgentConfig(defCoreConf("prg/dir", "app/dir", tmpJson(j)))
	actual := c.String()
	expected := fmt.Sprintf(`{coreConfig:%v, port:123, appKey:...}`, c.coreConfig.String())
	if actual != expected {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}

}

func TestString_appKey4letter(t *testing.T) {
	j := `
{
	"port": 123,
	"key": "9876"
}`
	c := loadAgentConfig(defCoreConf("prg/dir", "app/dir", tmpJson(j)))
	actual := c.String()
	expected := fmt.Sprintf(`{coreConfig:%v, port:123, appKey:987...}`, c.coreConfig.String())
	if actual != expected {
		t.Errorf("expect %v, but was %v.", expected, actual)
	}

}
