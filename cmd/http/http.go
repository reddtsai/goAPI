package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	_rootCmd = &cobra.Command{
		Use:   "api",
		Short: "http server",
	}
	_httpHandler http.Handler
)

func init() {
	port := 80
	_rootCmd.Flags().IntVar(&port, "port", port, "http listen port")
}

func runServer(cmd *cobra.Command, args []string) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cmd.Flag("port").Value.String()),
		Handler: _httpHandler,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server, listen and serve error : %s\n", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	err := server.Shutdown(cmd.Context())
	if err != nil {
		log.Printf("http server, shutdown error : %s\n", err.Error())
	}
}

func Execute(ctx context.Context, httpHandler http.Handler) error {
	_httpHandler = httpHandler
	_rootCmd.Run = runServer

	return _rootCmd.ExecuteContext(ctx)
}
