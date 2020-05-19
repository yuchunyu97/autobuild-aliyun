package auth

import (
	"io/ioutil"
	"os/user"
	"path"

	"gopkg.in/yaml.v2"
)

const configFile = ".iproxy.cred"

var homeDir string

// Credential 阿里云认证信息
type Credential struct {
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
}

func init() {
	user, err := user.Current()
	if err != nil {
		homeDir = "/tmp"
	}
	homeDir = user.HomeDir
}

// Save 保存为文件，或更新
func (c *Credential) Save() (err error) {
	byteArray, err := yaml.Marshal(c)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path.Join(homeDir, configFile), byteArray, 0644)
	if err != nil {
		return
	}

	return
}

// Get 从文件中读取
func (c *Credential) Get() (err error) {
	yamlFile, err := ioutil.ReadFile(path.Join(homeDir, configFile))
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return
	}

	return
}
