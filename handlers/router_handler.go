package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/naighes/imposter/cfg"
	"github.com/naighes/imposter/functions"
)

type RouterHandler struct {
	routes       []*route
	vars         map[string]interface{}
	storeHandler StoreHandler
}

type route struct {
	expression functions.Expression
	latency    time.Duration
	handler    http.Handler
}

func (router *RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.storeHandler != nil && router.storeHandler.ServeHTTP(w, r) {
		return
	}
	for _, route := range router.routes {
		// TODO: X-Forwarded-Host?
		ctx := &functions.EvaluationContext{Vars: router.vars, Req: r}
		a, err := route.expression.Evaluate(ctx)
		if err != nil {
			writeError(w, err)
			return
		}
		b, ok := a.(bool)
		if !ok {
			writeError(w, fmt.Errorf("rule_expression requires a 'bool' expression: found '%v' instead", reflect.TypeOf(a)))
			return
		}
		if b {
			if route.latency > 0 {
				time.Sleep(route.latency * time.Millisecond)
			}
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// TODO: not sure about just returning not found...
	http.NotFound(w, r)
}

func (router *RouterHandler) add(expression functions.Expression, latency time.Duration, h func(http.ResponseWriter, *http.Request)) {
	router.routes = append(router.routes, &route{expression, latency, http.HandlerFunc(h)})
}

func NewRouterHandler(config *cfg.Config, storeHandler StoreHandler) (*RouterHandler, error) {
	defs := config.Defs
	var vars map[string]interface{}
	if config.Vars == nil {
		vars = make(map[string]interface{})
	} else {
		vars = config.Vars
	}
	r := RouterHandler{}
	r.vars = vars
	r.storeHandler = storeHandler
	for _, def := range defs {
		rule, err := functions.ParseExpression(def.RuleExpression)
		if err != nil {
			return nil, err
		}
		f, err := HandleFunc(def.Response, vars)
		if err != nil {
			return nil, err
		}
		if def.Latency < 0 {
			return nil, fmt.Errorf("latency requires a value greater than zero")
		}
		r.add(rule, def.Latency, f)
	}
	return &r, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	w.WriteHeader(500)
	fmt.Fprintf(w, err.Error())
}

func HandleFunc(o interface{}, vars map[string]interface{}) (func(http.ResponseWriter, *http.Request), error) {
	var rsp cfg.MatchRsp
	err := mapstructure.Decode(o, &rsp)
	if err == nil {
		return matchRspHTTPHandler{content: &rsp, vars: vars}.handleFunc(functions.ParseExpression)
	}
	str, ok := o.(string)
	if ok {
		return funcHTTPHandler{content: str, vars: vars}.handleFunc(functions.ParseExpression)
	}
	return nil, fmt.Errorf("operation is not supported")
}
