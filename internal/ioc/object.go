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
	optional bool
	isFunc   bool
	function any
}

func NewObject(id string, name string, obj any) *Object {
	return &Object{id: id, name: name, obj: obj}
}

func NewObjectFromFunc(id string, name string, function any) *Object {
	return &Object{id: id, name: name, function: function, isFunc: true}
}

func (o *Object) SetInjected(injected bool) {
	o.injected = injected
}

func (o *Object) SetFuncRet(ret any) {
	o.obj = ret
}

func (o *Object) GetReflectType() reflect.Type {
	if o.isFunc {
		return reflect.TypeOf(o.function)
	}
	return reflect.TypeOf(o.obj).Elem()
}
func (o *Object) GetReflectValue() reflect.Value {
	if o.isFunc {
		return reflect.ValueOf(o.function)
	}
	return reflect.ValueOf(o.obj).Elem()
}

func (o *Object) GetValue() any {
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
