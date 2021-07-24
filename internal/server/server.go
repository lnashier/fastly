package server

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"time"

	gmux "github.com/gorilla/mux"
)

// Server serves the clients
type Server struct {
	*http.Server
	cfg *viper.Viper
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
	fmt.Println("server@Stop going to sleep for gracetime")
	time.Sleep(time.Second * time.Duration(s.cfg.GetInt("server.app.shutdown.gracetime")))
	fmt.Println("server@Stop wokeup after gracetime")
	err := s.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("server@Stop error %s\n", err.Error())
		return
	}
	fmt.Println("server@Stop server shutdown completed")
}

// New returns a new Server but doesn't start it.
// Call Start from outside. Call Stop to shut it down.
func New(cfg *viper.Viper) *Server {
	fmt.Println("server@New enter")
	defer fmt.Println("server@New exit")

	return &Server{
		Server: &http.Server{
			Addr: ":8080",
			Handler: (&mux{
				Router: gmux.NewRouter(),
				cfg:    cfg,
			}).init(),
		},
		cfg: cfg,
	}
}
