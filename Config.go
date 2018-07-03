package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Options *ConfigOptions `json:"options"`
	Defs    []*MatchDef    `json:"pattern_list"`
}

type ConfigOptions struct {
	Cors bool `json:"enable_cors"`
}

type MatchDef struct {
	Pattern  string      `json:"pattern"`
	Response interface{} `json:"response"`
}

type MatchRsp struct {
	Body       string                 `mapstructure:"body"`
	Headers    map[string]interface{} `mapstructure:"headers"`
	StatusCode int                    `mapstructure:"status_code"`
}

func ParseConfig(j []byte) (*Config, error) {
	var r *Config
	err := json.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rsp *MatchRsp) ParseBody() func() (string, error) {
	body := rsp.Body
	return func() (string, error) {
		name, arg, err := ParseFunc(body)
		if err != nil {
			return "", err
		}
		switch name {
		case "text":
			return arg, nil
		case "file":
			content, err := ioutil.ReadFile(arg)
			if err != nil {
				return "", err
			}
			return string(content), nil
		default:
			return "", fmt.Errorf("function '%s' is not supported", name)
		}
	}
}
