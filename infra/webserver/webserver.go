package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]Route
	WebServerPort string
}

func New(port string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]Route),
		WebServerPort: port,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc, method string) {
	var key = fmt.Sprintf("%s-%s", method, path)
	s.Handlers[key] = Route{
		Path:    path,
		Handler: handler,
		Method:  method,
	}
}

func (s *WebServer) AddMiddleware(middleware func(handler http.Handler) http.Handler) {
	s.Router.Use(middleware)
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	s.AddMiddleware(middleware.Recoverer)
	for _, handler := range s.Handlers {
		s.Router.MethodFunc(handler.Method, handler.Path, handler.Handler)
	}
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
