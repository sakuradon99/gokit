package ioc

type RegisterOptions struct {
	Name                string
	ImplementInterfaces []any
}

type RegisterOption func(o *RegisterOptions)
