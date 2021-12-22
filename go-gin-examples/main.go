package main

import (
	"cathub.me/go-web-examples/pkg/setting"
	"cathub.me/go-web-examples/timer"
	"cathub.me/go-web-examples/web"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	setting.Setup()
}

func main() {
	gin.SetMode(setting.Server.Env)
	r := web.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Server.Port),
		Handler:        r,
		ReadTimeout:    setting.Server.ReadTimeout,
		WriteTimeout:   setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil {
			log.Err(err).Msg("listen failed")
		}
	}()

	// Start timer task
	timer.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Server Shutdown")
	}
	log.Info().Msg("Server exiting")
}
