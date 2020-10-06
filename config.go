package war

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Directory     string   `yaml:"directory"`
	Match         []string `yaml:"match"`
	Exclude       []string `yaml:"exclude"`
	CommandString string   `yaml:"command"`
	Env           []string `yaml:"env"`
	Delay         int      `yaml:"delay"`
	Timeout       int      `yaml:"timeout"`
}

func ConfigFromFilename(filename string) (config, error) {
	var c config

	f, err := os.Open(filename)
	if err != nil {
		return c, fmt.Errorf("could not create config from file %s: %w", filename, err)
	}

	return ConfigFromFile(f)
}

func ConfigFromFile(f io.Reader) (config, error) {
	var c config

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return c, fmt.Errorf("could not create config from file: %w", err)
	}

	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return c, fmt.Errorf("could not create config from file: %w", err)
	}

	return c, nil
}
