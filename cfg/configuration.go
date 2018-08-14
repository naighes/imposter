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

// Config represents an imPOSTer configuration.
// A set of rule expressions can be defined dy Defs field.
type Config struct {
	Defs []*MatchDef            `json:"pattern_list" yaml:"pattern_list"`
	Vars map[string]interface{} `json:"vars" yaml:"vars"`
}

// MatchDef represents a single rule expression.
// The RuleExpression field wraps a boolean expression every incoming HTTP request is matched against.
// How a matching rule expression should be managed is defined by the Response object.
type MatchDef struct {
	RuleExpression string        `json:"rule_expression" yaml:"rule_expression"`
	Latency        time.Duration `json:"latency" yaml:"latency"`
	Response       interface{}   `json:"response" yaml:"response"`
}

// MatchRsp is the fully structured version of a Response object.
// Body represents the payload to be returned and it can be an expression as well.
// Headers is a collection of HTTP headers to be returned and each entry can be an expression.
// StatusCode represents the resulting HTTP status code and it MUST be an expression:
//		rsp := MatchRsp{Body: "some content", StatusCode: `${200}`}
type MatchRsp struct {
	Body       string                 `mapstructure:"body"`
	Headers    map[string]interface{} `mapstructure:"headers"`
	StatusCode string                 `mapstructure:"status_code"`
	Cookies    map[string]HTTPCookie  `mapstructure:"cookies"`
}

type HTTPCookie struct {
	Value   string `json:"value" yaml:"value"`
	Path    string `json:"path" yaml:"path"`
	Domain  string `json:"domain" yaml:"domain"`
	Expires string `json:"expires" yaml:"expires"`
}

func parseConfig(j []byte) (*Config, error) {
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

// ReadConfig takes a path as an input and parses its content to build the imPOSTer configuration.
func ReadConfig(configFile string) (*Config, error) {
	if configFile == "" {
		return &Config{}, nil
	}
	var err error
	var configPath string
	var rawConfig []byte
	var config *Config
	if configPath, err = filepath.Abs(configFile); err == nil {
		if rawConfig, err = ioutil.ReadFile(configPath); err == nil {
			if config, err = parseConfig(rawConfig); err == nil {
				return config, nil
			}
		}
	}
	return nil, err
}

// Validate method parses the current expression trying to catch potential evaluation errors.
// An empty array is returned whether no errors were found.
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
	_, err = rsp.ParseCookies(parse)
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

// ParseHeaders evaluates any HTTP header and returns an expression for each of them.
// It returns an error in case of evaluation failures.
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

func (rsp *MatchRsp) ParseCookies(parse functions.ExpressionParser) (map[string]map[string]functions.Expression, error) {
	if rsp.Cookies == nil {
		return map[string]map[string]functions.Expression{}, nil
	}
	cookies := make(map[string]map[string]functions.Expression)
	var errors error
	for k, c := range rsp.Cookies {
		cookie := make(map[string]functions.Expression)
		v, err := parse(c.Value)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			cookie["value"] = v
		}
		p, err := parse(c.Path)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			cookie["path"] = p
		}
		d, err := parse(c.Domain)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			cookie["domain"] = d
		}
		e, err := parse(c.Expires)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			cookie["expires"] = e
		}
		cookies[k] = cookie
	}
	if errors != nil {
		return nil, errors
	}
	return cookies, nil
}

// ParseStatusCode evaluates the resulting status code expression to an integer.
// It returns an error in case of evaluation failures.
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
