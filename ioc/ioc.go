package ioc

import "github.com/sakuradon99/gokit/internal/ioc"

var container ioc.Container = ioc.NewContainerImpl()

func Register(object any, opts ...ioc.RegisterOption) {
	err := container.Register(object, opts...)
	if err != nil {
		panic(err)
	}
}

func GetObject[T any](name string) *T {
	t := new(T)
	obj, err := container.GetObject(name, t)
	if err != nil {
		panic(err)
	}

	o, ok := obj.(*T)
	if !ok {
		panic("type assertion failed")
	}
	return o
}
