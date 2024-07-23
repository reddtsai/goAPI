package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "http server",
	Run:   runServer,
}

func init() {
	port := 80
	rootCmd.Flags().IntVar(&port, "port", port, "http listen port")
}

func runServer(cmd *cobra.Command, args []string) {
	mux := blockActionMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cmd.Flag("port").Value.String()),
		Handler: mux,
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

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
