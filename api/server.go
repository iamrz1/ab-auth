package api

import (
	"context"
	"fmt"
	"github.com/iamrz1/ab-auth/config"
	"github.com/iamrz1/ab-auth/service"
	"github.com/iamrz1/ab-auth/utils"
	rLog "github.com/iamrz1/rest-log"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/iamrz1/ab-auth/api/health"
)

func Start(cfg *config.AppConfig, svc *service.Config, logger rLog.Logger) (*http.Server, error) {

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	handler, err := SetupRouter(cfg, svc, logger)
	if err != nil {
		log.Println("cant setup router:", err)
		return nil, err
	}

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      handler,
	}

	go func() {
		log.Println("Staring server with address ", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Stopping server:", err)
			os.Exit(-1)
		}
	}()

	return srv, nil
}

func Stop(server *http.Server) error {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Http server couldn't shutdown gracefully", err)
		return err
	}

	log.Println("shutting down")
	return nil
}

func SetupRouter(cfg *config.AppConfig, svc *service.Config, logger rLog.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// enforce cors
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"*"},
		AllowOriginFunc:  verifyOrigin,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if cfg.Environment != utils.EnvProduction {
		r.Get("/doc/*", httpSwagger.Handler())
	}

	r.Mount("/", health.Router())
	r.Mount("/api/v1", V1Router(svc, logger))

	return r, nil
}

func verifyOrigin(r *http.Request, origin string) bool {
	log.Println("cors from ", origin)
	// todo: write a function to allow only valid origins
	return true
}
