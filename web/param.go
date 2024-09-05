package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

var validate = validator.New()

type paramBuilder interface {
	Build(ctx *gin.Context) (reflect.Value, error)
}

type defaultParamBuilder struct {
	rtp reflect.Type
}

func newDefaultParamBuilder(rtp reflect.Type) *defaultParamBuilder {
	return &defaultParamBuilder{rtp: rtp}
}

func (b *defaultParamBuilder) Build(_ *gin.Context) (reflect.Value, error) {
	return reflect.New(b.rtp).Elem(), nil
}

type contextParamBuilder struct {
}

func newContextParamBuilder() *contextParamBuilder {
	return &contextParamBuilder{}
}

func (b *contextParamBuilder) Build(ctx *gin.Context) (reflect.Value, error) {
	return reflect.ValueOf(ctx), nil
}

type structParamBuilder struct {
	rtp     reflect.Type
	binders []binder
}

func newStructParamBuilder(rtp reflect.Type, binders []binder) *structParamBuilder {
	return &structParamBuilder{rtp: rtp, binders: binders}
}

func (b *structParamBuilder) Build(ctx *gin.Context) (reflect.Value, error) {
	param := reflect.New(b.rtp).Elem()
	for _, binder := range b.binders {
		if err := binder.Bind(ctx, param); err != nil {
			return reflect.Value{}, err
		}
	}
	return param, nil
}
