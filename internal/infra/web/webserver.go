package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WebServerInterface interface {
	Start()
}

type RouteHandler struct {
	Path        string
	Method      string
	HandlerFunc http.HandlerFunc
}

type Middleware struct {
	Name    string
	Handler func(next http.Handler) http.Handler
}

type WebServer struct {
	Router        chi.Router
	Handlers      []RouteHandler
	Middlewares   []Middleware
	WebServerPort int
}

func NewWebServer(
	serverPort int,
	handlers []RouteHandler,
	middlewares []Middleware,
) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      handlers,
		Middlewares:   middlewares,
		WebServerPort: serverPort,
	}
}

func (s *WebServer) Start() {
	for _, m := range s.Middlewares {
		s.Router.Use(m.Handler)
	}

	for _, h := range s.Handlers {
		s.Router.MethodFunc(h.Method, h.Path, h.HandlerFunc)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", s.WebServerPort), s.Router)
}
