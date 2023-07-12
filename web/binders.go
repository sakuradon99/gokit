package web

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type binder interface {
	Bind(c *gin.Context, param reflect.Value) error
}

type pathBinder struct {
	field      int
	path       string
	defaultVal string
}

func (b *pathBinder) Bind(c *gin.Context, param reflect.Value) error {
	val := c.Param(b.path)
	if val == "" {
		val = b.defaultVal
	}
	err := setVal(param.Field(b.field), val)
	return err
}

type queryBinder struct {
	field      int
	query      string
	defaultVal string
}

func (b *queryBinder) Bind(c *gin.Context, param reflect.Value) error {
	val := c.Query(b.query)
	if val == "" {
		val = b.defaultVal
	}
	err := setVal(param.Field(b.field), val)
	return err
}

type formBinder struct {
	field      int
	form       string
	defaultVal string
}

func (b *formBinder) Bind(c *gin.Context, param reflect.Value) error {
	val := c.PostForm(b.form)
	if val == "" {
		val = b.defaultVal
	}
	err := setVal(param.Field(b.field), val)
	return err
}

type requestJSONBinder struct {
	field int
	body  reflect.Type
}

func (b *requestJSONBinder) Bind(c *gin.Context, param reflect.Value) error {
	body := reflect.New(b.body).Interface()
	err := c.ShouldBindJSON(body)
	if err != nil {
		return ErrInvalidParams
	}
	if validator, ok := body.(Validator); ok {
		err = validator.Validate()
		if err != nil {
			return ErrInvalidParams
		}
	}
	param.Field(b.field).Set(reflect.ValueOf(body).Elem())
	return nil
}

type fileBinder struct {
	field int
	file  string
}

func (f *fileBinder) Bind(c *gin.Context, param reflect.Value) error {
	file, err := c.FormFile(f.file)
	if err != nil {
		return ErrInvalidParams
	}

	param.Field(f.field).Set(reflect.ValueOf(file))
	return nil
}

type multipleFileBinder struct {
	field int
	files string
}

func (m *multipleFileBinder) Bind(c *gin.Context, param reflect.Value) error {
	form, err := c.MultipartForm()
	if err != nil {
		return ErrInvalidParams
	}

	param.Field(m.field).Set(reflect.ValueOf(form.File[m.files]))
	return nil
}
