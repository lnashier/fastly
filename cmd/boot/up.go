package boot

import (
	"fmt"
	"github.com/fastly/internal/server"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// Up is to bootup server
// It is a blocking function
func Up(cfg *viper.Viper) {
	fmt.Println("boot@Up enter")
	defer fmt.Println("boot@Up exit")

	srv := server.New(cfg)

	go func(srv *server.Server) {
		// Logging: Info Level

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc)

	Loop:
		for {
			sig := <-sigc
			// kill -SIGXXX <pid>
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT:
				fmt.Printf("boot@Up signal (%v) received\n", sig.String())
				break Loop
			default:
				fmt.Printf("boot@Up signal (%v) ignored\n", sig.String())
			}
		}
		fmt.Println("boot@Up stopping server")
		srv.Stop()
	}(srv)

	if err := srv.Start(); err != nil {
		fmt.Printf("boot@Up server start error %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}