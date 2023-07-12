package web

type Options func(i *Interceptor)

func WithTplSuffix(suffix string) Options {
	return func(i *Interceptor) {
		i.tplSuffix = suffix
	}
}

func WithViewIntercept(f viewIntercept) Options {
	return func(i *Interceptor) {
		i.viewIntercept = f
	}
}
