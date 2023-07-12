package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var validate = validator.New()

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

	fName := runtime.FuncForPC(fv.Pointer()).Name()
	if !validateHandlerFunc(ft) {
		panic(fmt.Errorf("wrong handler func: %s", fName))
	}

	ftNumIn := ft.NumIn()

	successHTTPCode := http.StatusOK
	var binders []binder
	var err error
	if ft.NumIn() == 2 {
		p := reflect.New(ft.In(1))
		pe := p.Elem()
		for i := 0; i < pe.NumField(); i++ {
			tag := pe.Type().Field(i).Tag
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
						body:  pe.Field(i).Type(),
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
			} else if code := tag.Get("success"); code != "" {
				successHTTPCode, err = strconv.Atoi(code)
				if err != nil || successHTTPCode < 100 || successHTTPCode > 999 {
					panic(fmt.Errorf("invalid http error code: %s in func %s", code, fName))
				}
			}
		}
	}

	return func(c *gin.Context) {
		incomes := make([]reflect.Value, ftNumIn)
		incomes[0] = reflect.ValueOf(c)

		if ftNumIn == 2 {
			param := reflect.New(ft.In(1))
			paramElem := param.Elem()
			for _, b := range binders {
				err := b.Bind(c, paramElem)
				if err != nil {
					handleError(c, err)
					return
				}
			}

			if validate.Struct(paramElem.Interface()) != nil {
				handleError(c, ErrInvalidParams)
				return
			}
			incomes[1] = paramElem
		}

		outcomes := fv.Call(incomes)

		// if handler does not return parameters, return client success
		if len(outcomes) == 0 {
			c.Status(successHTTPCode)
			return
		}

		response := outcomes[0]
		errResponse := outcomes[len(outcomes)-1]

		if err, ok := errResponse.Interface().(error); ok {
			handleError(c, err)
			return
		}

		it.handleResponse(c, response, successHTTPCode)
	}
}

func validateHandlerFunc(ft reflect.Type) bool {
	if ft.Kind() != reflect.Func || ft.NumIn() < 1 || ft.NumIn() > 2 || ft.NumOut() > 2 {
		return false
	}
	// the first income parameter must be context
	if ft.In(0).Name() != "Context" {
		return false
	}
	// if handler return 2 parameters, the second parameter must be error
	if ft.NumOut() == 2 && ft.Out(1).Name() != "error" {
		return false
	}

	return true
}

func handleError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func (it *Interceptor) handleResponse(c *gin.Context, response reflect.Value, statusCode int) {
	if (response.Kind() == reflect.Ptr || response.Kind() == reflect.Interface) && response.IsNil() {
		c.Status(statusCode)
		return
	}
	if response.Kind() == reflect.Slice && response.IsNil() {
		c.JSON(statusCode, make([]any, 0))
		return
	}
	if v, ok := response.Interface().(View); ok {
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
	if v, ok := response.Interface().(File); ok {
		c.File(v.Path)
		return
	}

	c.JSON(statusCode, response.Interface())
}
