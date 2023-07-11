package ioc

import "github.com/sakuradon99/gokit/internal/ioc"

func Name(name string) ioc.RegisterOption {
	return func(o *ioc.RegisterOptions) {
		o.Name = name
	}
}

func Implement(inf any) ioc.RegisterOption {
	return func(o *ioc.RegisterOptions) {
		o.ImplementInterfaces = append(o.ImplementInterfaces, inf)
	}
}
