package ioc

type RegisterOptions struct {
	Name                string
	Optional            bool
	ImplementInterfaces []any
}

type RegisterOption func(o *RegisterOptions)

func Name(name string) RegisterOption {
	return func(o *RegisterOptions) {
		o.Name = name
	}
}

func Implement(inf any) RegisterOption {
	return func(o *RegisterOptions) {
		o.ImplementInterfaces = append(o.ImplementInterfaces, inf)
	}
}

func Optional() RegisterOption {
	return func(o *RegisterOptions) {
		o.Optional = true
	}
}
