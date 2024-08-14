package web

type Handler interface {
	Base() string
	Routes() []Route
}
