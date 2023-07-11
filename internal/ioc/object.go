package ioc

import (
	"fmt"
	"reflect"
)

type Object struct {
	id       string
	name     string
	obj      any
	injected bool
	isFunc   bool
	funcRet  any
}

func NewObject(id string, name string, obj any, isFunc bool) *Object {
	return &Object{id: id, name: name, obj: obj, isFunc: isFunc}
}

func (o *Object) SetInjected(injected bool) {
	o.injected = injected
}

func (o *Object) SetFuncRet(ret any) {
	o.funcRet = ret
}

func (o *Object) GetReflectType() reflect.Type {
	if o.isFunc {
		return reflect.TypeOf(o.obj)
	}
	return reflect.TypeOf(o.obj).Elem()
}
func (o *Object) GetReflectValue() reflect.Value {
	if o.isFunc {
		return reflect.ValueOf(o.obj)
	}
	return reflect.ValueOf(o.obj).Elem()
}

func (o *Object) GetValue() any {
	if o.isFunc {
		return o.funcRet
	}
	return o.obj
}

type ObjectPool struct {
	objects map[string]*Object
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		objects: make(map[string]*Object),
	}
}

func (p *ObjectPool) Add(object *Object) error {
	if _, ok := p.objects[object.id]; ok {
		return fmt.Errorf("object with id %s already exists", object.id)
	}
	p.objects[object.id] = object
	return nil
}

func (p *ObjectPool) List() []*Object {
	var objects []*Object
	for _, obj := range p.objects {
		objects = append(objects, obj)
	}
	return objects
}

func (p *ObjectPool) Get(id string) (*Object, bool) {
	obj, ok := p.objects[id]
	return obj, ok
}

func genObjectID(pkgPath string, name string, alisa string) string {
	return fmt.Sprintf("%s.%s-%s", pkgPath, name, alisa)
}
