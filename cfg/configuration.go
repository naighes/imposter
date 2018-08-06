package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
	"github.com/naighes/imposter/functions"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Defs []*MatchDef            `json:"pattern_list" yaml:"pattern_list"`
	Vars map[string]interface{} `json:"vars" yaml:"vars"`
}

type MatchDef struct {
	RuleExpression string        `json:"rule_expression" yaml:"rule_expression"`
	Latency        time.Duration `json:"latency" yaml:"latency"`
	Response       interface{}   `json:"response" yaml:"response"`
}

type MatchRsp struct {
	Body       string                 `mapstructure:"body"`
	Headers    map[string]interface{} `mapstructure:"headers"`
	StatusCode string                 `mapstructure:"status_code"`
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

func ReadConfig(configFile string) (*Config, error) {
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

func (def *MatchDef) Validate(parse functions.ExpressionParser, vars map[string]interface{}) []string {
	var r []string
	if err := validateRuleExpression(def.RuleExpression, vars); err != nil {
		r = append(r, fmt.Sprintf("%v", err))
	}
	var rsp MatchRsp
	err := mapstructure.Decode(def.Response, &rsp)
	if err == nil {
		rr := rsp.validate(parse, vars)
		r = append(r, rr...)
	} else {
		body, _ := def.Response.(string)
		err := validateComputedBody(body, vars)
		if err != nil {
			r = append(r, fmt.Sprintf("%v", err))
		}
	}
	return r
}

func (rsp *MatchRsp) validate(parse functions.ExpressionParser, vars map[string]interface{}) []string {
	var r []string
	_, err := validateEvaluation(rsp.Body, vars)
	if err != nil {
		r = append(r, fmt.Sprintf("%v", err))
	}
	_, err = rsp.ParseHeaders(parse)
	if err != nil {
		if errors, ok := err.(*multierror.Error); ok {
			for err := range errors.Errors {
				r = append(r, fmt.Sprintf("%v", err))
			}
		}
	}
	err = validateStatusCode(rsp.StatusCode, vars)
	if err != nil {
		r = append(r, fmt.Sprintf("%v", err))
	}
	return r
}

func validateStatusCode(expression string, vars map[string]interface{}) error {
	if expression == "" {
		return nil
	}
	a, err := validateEvaluation(expression, vars)
	if err != nil {
		return err
	}
	if _, ok := a.(int); !ok {
		return fmt.Errorf("expected an 'int' value for status code; got '%v' instead", reflect.TypeOf(a))
	}
	return nil
}

func validateRuleExpression(expression string, vars map[string]interface{}) error {
	e, err := validateEvaluation(expression, vars)
	if err != nil {
		return err
	}
	_, ok := e.(bool)
	if !ok {
		return fmt.Errorf("evaluation error: expected 'bool' for any rule expression; got '%v' instead", reflect.TypeOf(e))
	}
	return nil
}

func validateComputedBody(expression string, vars map[string]interface{}) error {
	e, err := validateEvaluation(expression, vars)
	if err != nil {
		return err
	}
	_, ok := e.(*functions.HTTPRsp)
	if !ok {
		return fmt.Errorf("evaluation error: expected 'HTTPRsp' for a computed version of response object; got '%v' instead", reflect.TypeOf(e))
	}
	return nil
}

func validateEvaluation(expression string, vars map[string]interface{}) (interface{}, error) {
	e, err := functions.ParseExpression(expression)
	if err != nil {
		return nil, err
	}
	ctx := &functions.EvaluationContext{Vars: vars, Req: &http.Request{Header: http.Header{}}}
	a, err := e.Test(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (rsp *MatchRsp) ParseHeaders(parse functions.ExpressionParser) (map[string]functions.Expression, error) {
	headers := make(map[string]functions.Expression)
	var errors error
	if rsp.Headers != nil {
		for k, v := range rsp.Headers {
			header, ok := v.(string)
			if !ok {
				errors = multierror.Append(errors, fmt.Errorf("expected a value of type 'string'; got '%v' instead", reflect.TypeOf(v)))
				continue
			}
			he, err := parse(header)
			if err != nil {
				errors = multierror.Append(errors, err)
				continue
			}
			headers[k] = he
		}
	}
	if errors != nil {
		return nil, errors
	}
	return headers, nil
}

func (rsp *MatchRsp) ParseStatusCode(parse functions.ExpressionParser) (functions.Expression, error) {
	if rsp.StatusCode == "" {
		return parse("${200}")
	}
	e, err := parse(rsp.StatusCode)
	if err != nil {
		return nil, err
	}
	return e, nil
}
