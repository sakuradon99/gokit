package crud

type queryOptions struct {
	PreloadAssociations bool
}

type QueryOption func(*queryOptions)

func PreloadAssociations() QueryOption {
	return func(o *queryOptions) {
		o.PreloadAssociations = true
	}
}

func buildQueryOptions(opts ...QueryOption) queryOptions {
	option := &queryOptions{}
	for _, opt := range opts {
		opt(option)
	}
	return *option
}

type removeOptions struct {
	RemoveAssociations bool
}

type RemoveOption func(*removeOptions)

func RemoveAssociations() RemoveOption {
	return func(o *removeOptions) {
		o.RemoveAssociations = true
	}
}

func buildRemoveOptions(opts ...RemoveOption) removeOptions {
	option := &removeOptions{}
	for _, opt := range opts {
		opt(option)
	}
	return *option
}
