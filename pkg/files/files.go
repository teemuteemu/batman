package files

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Header map[string]string

type Request struct {
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	URL    string `yaml:"url"`
	Header Header `yaml:"header"`
	Body   string `yaml:"body"`
}

type Collection struct {
	Version     uint      `yaml:"version"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Requests    []Request `yaml:"requests"`
}

type Step struct {
	Name    string  `yaml:"name"`
	Request string  `yaml:"request"`
	Before  *string `yaml:"before"`
	After   *string `yaml:"after"`
	Output  bool    `yaml:"output"`
}

type Script struct {
	Version     uint   `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Collection  string `yaml:"collection"`
	Steps       []Step `yaml:"steps"`
}

func UnmarshalCollection(collectionStr []byte) (Collection, error) {
	collection := Collection{}

	err := yaml.Unmarshal(collectionStr, &collection)
	if err != nil {
		return collection, err
	}

	return collection, nil
}

func (c *Collection) FindRequest(requestName string) (*Request, error) {
	for _, request := range c.Requests {
		if request.Name == requestName {
			return &request, nil
		}
	}

	return nil, fmt.Errorf(`Request not found: "%s"`, requestName)
}

func GetCollection(collectionFile string) (*Collection, error) {
	collectionStr, err := ioutil.ReadFile(collectionFile)
	if err != nil {
		return nil, err
	}

	c, err := UnmarshalCollection(collectionStr)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func UnmarshalScript(flowStr []byte) (Script, error) {
	script := Script{}

	err := yaml.Unmarshal(flowStr, &script)
	if err != nil {
		return script, err
	}

	return script, nil
}

func GetScript(scriptFile string) (*Script, error) {
	scriptStr, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		return nil, err
	}

	s, err := UnmarshalScript(scriptStr)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
