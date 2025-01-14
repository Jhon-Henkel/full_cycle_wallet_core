package webserver

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type Webserver struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebserverPort string
}

func NewWebserver(router chi.Router, handlers map[string]http.HandlerFunc, webserverPort string) *Webserver {
	return &Webserver{
		Router:        router,
		Handlers:      handlers,
		WebserverPort: webserverPort,
	}
}

func (s *Webserver) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *Webserver) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Post(path, handler)
	}
	http.ListenAndServe(s.WebserverPort, s.Router)
}
