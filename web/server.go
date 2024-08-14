package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sakuradon99/ioc"
	"path"
	"sort"
	"strings"
)

var _ = ioc.Register[Server]()

type Server struct {
	host string `value:"web.host;optional"`
	port string `value:"web.port;optional"`

	handlers            []Handler                  `inject:"r:.*"`
	middlewares         []Middleware               `inject:"r:.*"`
	customEngineConfigs []ServerCustomEngineConfig `inject:"r:.*"`
}

func (s *Server) Init() error {
	if s.host == "" {
		s.host = "0.0.0.0"
	}
	if s.port == "" {
		s.port = "8080"
	}
	return nil
}

func (s *Server) Run() error {
	wi := NewInterceptor()
	server := gin.Default()

	for _, config := range s.customEngineConfigs {
		err := config.CustomEngine(server)
		if err != nil {
			return err
		}
	}

	var routePaths []string
	for _, handler := range s.handlers {
		rootPath := handler.Base()
		if !strings.HasPrefix(rootPath, "/") {
			rootPath = "/" + rootPath
		}

		for _, route := range handler.Routes() {
			var routePath string
			if strings.HasPrefix(route.Path, ".") {
				routePath = rootPath + route.Path
			} else {
				routePath = path.Join(rootPath, route.Path)
			}

			ginHandlers := s.applyMiddlewares(routePath)
			ginHandlers = append(ginHandlers, wi.Intercept(route.Func))

			routePaths = append(routePaths, routePath)

			server.Handle(route.Method, routePath, ginHandlers...)
		}
	}

	sort.Slice(routePaths, func(i, j int) bool {
		return routePaths[i] < routePaths[j]
	})

	return server.Run(fmt.Sprintf("%s:%s", s.host, s.port))
}

func (s *Server) applyMiddlewares(path string) []gin.HandlerFunc {
	var handlers []middlewareHandlerWithOrder
	for _, middleware := range s.middlewares {
		order := middleware.Register(path)
		if order < 0 {
			continue
		}
		handlers = append(handlers, middlewareHandlerWithOrder{
			fn:    middleware.Handle,
			order: order,
		})
	}

	sort.Slice(handlers, func(i, j int) bool {
		return handlers[i].order < handlers[j].order
	})

	var ginHandlers []gin.HandlerFunc
	for _, handler := range handlers {
		ginHandlers = append(ginHandlers, handler.fn)
	}

	return ginHandlers
}

func Run() error {
	server, err := ioc.GetObject[Server]("")
	if err != nil {
		return err
	}

	return server.Run()
}
