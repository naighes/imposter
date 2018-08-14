package functions

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type rndStringFunction struct {
	length Expression
}

func newRndStringFunction(args []Expression) (Expression, error) {
	if l := len(args); l != 1 {
		return nil, fmt.Errorf("function 'rnd_string' is expecting one argument of type 'int'; found %d argument(s) instead", l)
	}
	r := rndStringFunction{length: args[0]}
	return r, nil
}

func (f rndStringFunction) evaluate(g func(Expression) (interface{}, error), req *http.Request) (interface{}, error) {
	a, err := g(f.length)
	if err != nil {
		return "", err
	}
	b, ok := a.(int)
	if !ok {
		return "", fmt.Errorf("evaluation error: cannot convert value '%v' to 'int'", a)
	}
	if b < 0 {
		return "", fmt.Errorf("evaluation error: expected a positive integer value; got '%v' instead", b)
	}
	return randStringBytesMaskImprSrc(b), nil
}

func (f rndStringFunction) Evaluate(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Evaluate(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}

func (f rndStringFunction) Test(ctx *EvaluationContext) (interface{}, error) {
	g := func(ctx *EvaluationContext) func(Expression) (interface{}, error) {
		return func(expression Expression) (interface{}, error) {
			return expression.Test(ctx)
		}
	}(ctx)
	return f.evaluate(g, ctx.Req)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
