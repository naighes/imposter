package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/naighes/imposter/functions"
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

func readConfig(configFile string) (*Config, error) {
	var err error
	var configPath string
	var rawConfig []byte
	var config *Config
	if configPath, err = filepath.Abs(configFile); err == nil {
		if rawConfig, err = ioutil.ReadFile(configPath); err == nil {
			if config, err = ParseConfig(rawConfig); err == nil {
				return config, nil
			}
		}
	}
	return &Config{}, err
}

func (def *MatchDef) validate(vars map[string]interface{}) []string {
	var r []string
	if err := validateRuleExpression(def.RuleExpression, vars); err != nil {
		r = append(r, fmt.Sprintf("%v", err))
	}

	var rsp MatchRsp
	err := mapstructure.Decode(def.Response, &rsp)
	if err == nil {
		rr := rsp.validate(vars)
		r = append(r, rr...)
	}

	return r
}

func (rsp *MatchRsp) validate(vars map[string]interface{}) []string {
	var r []string
	_, err := validateEvaluation(rsp.Body, vars)
	if err != nil {
		r = append(r, fmt.Sprintf("%v", err))
	}
	if rsp.Headers != nil {
		for _, v := range rsp.Headers {
			header, ok := v.(string)
			if !ok {
				r = append(r, fmt.Sprintf("expected a value of type 'string'; got '%s' instead", reflect.TypeOf(v)))
			} else {
				_, err := validateEvaluation(header, vars)
				if err != nil {
					r = append(r, fmt.Sprintf("%v", err))
				}
			}
		}
	}
	return r
}

func validateRuleExpression(expression string, vars map[string]interface{}) error {
	e, err := validateEvaluation(expression, vars)
	if err != nil {
		return err
	}
	_, ok := e.(bool)
	if !ok {
		fmt.Errorf("evaluation error: expected 'bool' for any rule expression; got '%v' instead", reflect.TypeOf(e))
	}
	return nil
}

func validateEvaluation(expression string, vars map[string]interface{}) (interface{}, error) {
	e, err := functions.ParseExpression(expression)
	if err != nil {
		return nil, err
	}
	a, err := e.Test(vars, &http.Request{Header: http.Header{}})
	if err != nil {
		return nil, err
	}
	return a, nil
}
