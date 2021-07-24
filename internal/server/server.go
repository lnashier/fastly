package server

import (
	"context"
	"fmt"
	"net/http"

	gmux "github.com/gorilla/mux"
)

// Server serves the clients
type Server struct {
	*http.Server
}

// Start is to start the server
func (s *Server) Start() error {
	fmt.Println("server@Start enter")
	defer fmt.Println("server@Start exit")
	return s.ListenAndServe()
}

// Stop is to stop the server
func (s *Server) Stop() {
	fmt.Println("server@Stop enter")
	defer fmt.Println("server@Stop exit")

	err := s.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("server@Stop error %s\n", err.Error())
		return
	}
}

// New returns a new Server but doesn't start it.
// Call Start from outside. Call Stop to shut it down.
func New() *Server {
	fmt.Println("server@New enter")
	defer fmt.Println("server@New exit")

	return &Server{
		Server: &http.Server{
			Addr: ":8080",
			Handler: (&mux{
				Router: gmux.NewRouter(),
			}).init(),
		},
	}
}
