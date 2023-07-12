package web

import "context"

type View struct {
	Tpl  string
	Data any
}

type viewIntercept = func(ctx context.Context, v View) View
