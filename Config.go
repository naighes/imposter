package main

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Options *ConfigOptions         `json:"options" yaml:"options"`
	Defs    []*MatchDef            `json:"pattern_list" yaml:"pattern_list"`
	Vars    map[string]interface{} `json:"vars" yaml:"vars"`
}

type ConfigOptions struct {
	Cors bool `json:"enable_cors" yaml:"enable_cors"`
}

type MatchDef struct {
	RuleExpression string        `json:"rule_expression" yaml:"rule_expression"`
	Method         string        `json:"method" yaml:"method"`
	Latency        time.Duration `json:"latency" yaml:"latency"`
	Response       interface{}   `json:"response" yaml:"response"`
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
		return parseYaml(j)
	}
	return r, nil
}

func parseYaml(j []byte) (*Config, error) {
	var r *Config
	err := yaml.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
