package web

import "net/http"

type Route struct {
	Path   string
	Method string
	Func   any
}

func Get(path string, f any) Route {
	return Route{
		Path:   path,
		Method: http.MethodGet,
		Func:   f,
	}
}

func Post(path string, f any) Route {
	return Route{
		Path:   path,
		Method: http.MethodPost,
		Func:   f,
	}
}

func Put(path string, f any) Route {
	return Route{
		Path:   path,
		Method: http.MethodPut,
		Func:   f,
	}
}

func Delete(path string, f any) Route {
	return Route{
		Path:   path,
		Method: http.MethodDelete,
		Func:   f,
	}
}

func Routes(routes ...Route) []Route {
	return routes
}
