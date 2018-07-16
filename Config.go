package main

import (
	"encoding/json"
	"time"
)

type Config struct {
	Options *ConfigOptions         `json:"options"`
	Defs    []*MatchDef            `json:"pattern_list"`
	Vars    map[string]interface{} `json:"vars"`
}

type ConfigOptions struct {
	Cors bool `json:"enable_cors"`
}

type MatchDef struct {
	Pattern  string        `json:"pattern"`
	Method   string        `json:"method"`
	Latency  time.Duration `json:"latency"`
	Response interface{}   `json:"response"`
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
