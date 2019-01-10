package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application",
	Long:  "Serve the application",
	Run: func(cmd *cobra.Command, args []string) {
		wait := time.Second * 15

		r := mux.NewRouter()

		srv := &http.Server{
			Addr:         "0.0.0.0:8080",
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		srv.Shutdown(ctx)

		log.Println("shutting down")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
