package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sakuradon99/gokit/opt"
	"net/http"
	"reflect"
	"strings"
)

var (
	typeGinContext = reflect.TypeOf((*gin.Context)(nil))
	typeContext    = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeError      = reflect.TypeOf((*error)(nil)).Elem()
)

type Interceptor struct {
	tplSuffix     string
	viewIntercept viewIntercept
}

func NewInterceptor(options ...Options) *Interceptor {
	i := &Interceptor{}
	for _, option := range options {
		option(i)
	}
	return i
}

func (it *Interceptor) Intercept(f any) gin.HandlerFunc {
	ft := reflect.TypeOf(f)
	fv := reflect.ValueOf(f)
	ftNumIn := ft.NumIn()

	var paramBuilders []paramBuilder
	for i := 0; i < ftNumIn; i++ {
		field := ft.In(i)
		if field.Kind() == reflect.Struct {
			paramBuilders = append(paramBuilders, it.buildStructParamBuilder(field))
			continue
		}
		if field == typeGinContext || field.Implements(typeContext) {
			paramBuilders = append(paramBuilders, newContextParamBuilder())
			continue
		}
		paramBuilders = append(paramBuilders, newDefaultParamBuilder(field))
	}

	return func(c *gin.Context) {
		incomes := make([]reflect.Value, 0, ftNumIn)
		for _, builder := range paramBuilders {
			param, err := builder.Build(c)
			if err != nil {
				it.handleError(c, err)
				return
			}
			incomes = append(incomes, param)
		}

		outcomes := fv.Call(incomes)

		resp := opt.Empty[reflect.Value]()
		var err error
		for _, outcome := range outcomes {
			if outcome.Type().Implements(typeError) {
				e, ok := outcome.Interface().(error)
				if ok && e != nil {
					err = e
				}
			} else {
				resp = opt.Of(outcome)
			}
		}
		if err != nil {
			it.handleError(c, err)
			return
		}

		it.handleResponse(c, resp)
	}
}

func (it *Interceptor) buildStructParamBuilder(rtp reflect.Type) *structParamBuilder {
	var binders []binder
	for i := 0; i < rtp.NumField(); i++ {
		field := rtp.Field(i)
		tag := field.Tag
		defaultVal := tag.Get("default")
		if path := tag.Get("path"); path != "" {
			binders = append(binders, &pathBinder{
				field:      i,
				path:       path,
				defaultVal: defaultVal,
			})
		} else if query := tag.Get("query"); query != "" {
			binders = append(binders, &queryBinder{
				field:      i,
				query:      query,
				defaultVal: defaultVal,
			})
		} else if form := tag.Get("form"); form != "" {
			binders = append(binders, &formBinder{
				field:      i,
				form:       form,
				defaultVal: defaultVal,
			})
		} else if request := tag.Get("request"); request != "" {
			if request == "json" {
				binders = append(binders, &requestJSONBinder{
					field: i,
					body:  field.Type,
				})
			}
		} else if file := tag.Get("file"); file != "" {
			binders = append(binders, &fileBinder{
				field: i,
				file:  file,
			})
		} else if files := tag.Get("files"); files != "" {
			binders = append(binders, &multipleFileBinder{
				field: i,
				files: files,
			})
		}
	}

	return newStructParamBuilder(rtp, binders)
}

func (it *Interceptor) handleError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func (it *Interceptor) handleResponse(c *gin.Context, resp opt.Optional[reflect.Value]) {
	if !resp.Exists() {
		c.Status(http.StatusNoContent)
		return
	}
	if (resp.Get().Kind() == reflect.Ptr || resp.Get().Kind() == reflect.Interface) && resp.Get().IsNil() {
		c.Status(http.StatusNoContent)
		return
	}
	if resp.Get().Kind() == reflect.Slice && resp.Get().IsNil() {
		c.JSON(http.StatusOK, make([]any, 0))
		return
	}
	// TODO refactor the view handler
	if v, ok := resp.Get().Interface().(View); ok {
		if it.tplSuffix != "" && !strings.HasSuffix(v.Tpl, it.tplSuffix) {
			v.Tpl = v.Tpl + it.tplSuffix
		}
		if v.Data == nil {
			v.Data = gin.H{}
		}
		if it.viewIntercept != nil {
			v = it.viewIntercept(c, v)
		}

		c.HTML(http.StatusOK, v.Tpl, v.Data)
		return
	}
	if v, ok := resp.Get().Interface().(File); ok {
		c.File(v.Path)
		return
	}

	c.JSON(http.StatusOK, resp.Get().Interface())
}
