package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

func PullYml(source string, conf interface{}) error {
	file, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		return err
	}
	return nil
}
