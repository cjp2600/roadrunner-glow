package glow

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"github.com/yalp/jsonpath"
	"gopkg.in/workanator/go-floc.v2"
	"net/http"
	"regexp"
	"strings"
)

// VariablePrefix variable token prefix
const VariablePrefix = "$"

// RequestTranslator request structure
type RequestTranslator struct {
	Job     *Job
	Options *Option
	Ctx     context.Context
	client  *resty.Client
	logger  zerolog.Logger
}

// NewRequestTranslator request constructor
func NewRequestTranslator(job *Job, logger zerolog.Logger, options *Option) *RequestTranslator {
	return &RequestTranslator{Job: job, client: resty.New(), logger: logger, Options: options}
}

// setType set variable type
func (r *RequestTranslator) setType(ctx floc.Context, n, t string, val interface{}) {
	switch strings.ToLower(t) {
	case "string":
		ctx.AddValue(VariablePrefix+n, val.(string))
	case "int":
		ctx.AddValue(VariablePrefix+n, val.(int))
	case "int32":
		ctx.AddValue(VariablePrefix+n, val.(int32))
	case "int64":
		ctx.AddValue(VariablePrefix+n, val.(int64))
	case "float64":
		ctx.AddValue(VariablePrefix+n, val.(float64))
	}
}

// FindVarByJPath find variable in json
func (r *RequestTranslator) FindVarByJPath(js string, jPath string) (interface{}, error) {
	var store interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(js), &store)
	if err != nil {
		return nil, err
	}
	val, err := jsonpath.Read(store, jPath)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// FlocExecute execute
func (r *RequestTranslator) FlocExecute() func(ctx floc.Context, ctrl floc.Control) error {
	return func(ctx floc.Context, ctrl floc.Control) error {
		resp, err := r.Execute(ctx)
		if err != nil {
			return err
		}
		ctx.AddValue(r.Job.Id, resp.String())
		if r.Job.Var != nil {
			for _, v := range r.Job.Var {
				if len(v.JPath) > 0 {
					val, err := r.FindVarByJPath(resp.String(), v.JPath)
					if err != nil {
						return err
					}
					r.setType(ctx, v.Name, v.Type, val)
				}
			}
		}
		return nil
	}
}

// FindVariables find variable
func (r *RequestTranslator) FindVariables(ctx floc.Context, s string) string {
	reg, _ := regexp.Compile(`[$][a-zA-Z+]*`)
	v := reg.FindAllString(s, -1)
	if len(v) > 0 {
		for _, item := range v {
			if v, ok := ctx.Value(item).(string); ok {
				s = strings.Replace(s, item, v, -1)
			}
		}
	}
	return s
}

// Execute
func (r *RequestTranslator) Execute(ctx floc.Context) (*resty.Response, error) {
	var err error
	var response *resty.Response
	req := r.client.R()
	req.EnableTrace()

	url := r.Job.Url
	if strings.Contains(url, VariablePrefix) {
		url = r.FindVariables(ctx, url)
	}

	if r.Job.Body != nil {
		for _, value := range r.Job.Body {
			if v, ok := value.(string); ok {
				if strings.Contains(v, VariablePrefix) {
					v = r.FindVariables(ctx, v)
				}
			}
		}
		req.SetBody(r.Job.Body)
	}

	if r.Job.Header != nil {
		var headers = make(map[string]string)
		for key, value := range r.Job.Header {
			if sv, ok := value.(string); ok {
				if strings.Contains(sv, VariablePrefix) {
					headers[key] = r.FindVariables(ctx, sv)
				}
			}
		}
		req.SetHeaders(headers)
	}

	switch strings.ToUpper(r.Job.Method) {
	case http.MethodGet:
		response, err = req.Get(url)
	case http.MethodPost:
		response, err = req.Post(url)
	case http.MethodPut:
		response, err = req.Put(url)
	case http.MethodDelete:
		response, err = req.Delete(url)
	}

	if err != nil {
		return nil, err
	}

	if response == nil {
		r.logger.Debug().Str("type", "request").Err(errors.New("empty response"))
	}
	if r.Options.Debug {
		r.logger.Debug().Str("type", "request").
			Str("method", r.Job.Method).
			Str("url", r.Job.Url).
			Dur("connTime", req.TraceInfo().ConnTime).
			Dur("serverTime", req.TraceInfo().ServerTime).
			Dur("responseTime", req.TraceInfo().ResponseTime).
			Dur("totalTime", req.TraceInfo().TotalTime).
			Msgf(string(response.Body()))
	}
	return response, nil
}
