package sftpsync

import (
	"encoding/json"

	"io/ioutil"
	"os"
	"path/filepath"
	_ "strings"
)

const (
	conf = "config.json"
)

type Config struct {
	Host       string
	User       string
	Password   string
	Method     string
	LocalPath  string
	ServerPath string
}

func NewConfig(confFile string) (config *Config, err error) {
	//baseDir := filepath.Dir(args[0])
	return ParseConfig(confFile)
}

func ParseConfig(confFile string) (*Config, error) {
	if ok, _ := isFileExists(confFile); !ok {
		confFile = filepath.Join(baseDir, conf)
	}
	config, err := readConfigFile(confFile)
	if config.ServerPath == "" {
		config.ServerPath = "/"
	}
	if config.LocalPath == "" {
		config.LocalPath = filepath.Join(baseDir, "data")
	}
	return config, err
}
func readConfigFile(confFile string) (config *Config, err error) {
	config = new(Config)
	f, err := os.Open(confFile)

	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(d, config); err != nil {
		return nil, err
	}

	//		if err := json.NewDecoder(f).Decode(config); err != nil {
	//			fmt.Println(err)
	//			return nil, err
	//		}

	return config, nil
}
