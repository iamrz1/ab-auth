package cmd

import (
	"context"
	"github.com/iamrz1/ab-auth/api"
	"github.com/iamrz1/ab-auth/config"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	infraMongo "github.com/iamrz1/ab-auth/infra/mongo"
	"github.com/iamrz1/ab-auth/service"
	rLog "github.com/iamrz1/rest-log"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SrvCmd is the serve sub command to start the api server
var SrvCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves the api server",
	RunE:  serve,
}

var db *infraMongo.Mongo
var cache *infraCache.Redis

func serve(cmd *cobra.Command, args []string) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server, err := StartServer()
	if err != nil {
		return err
	}

	err = StopServer(server)
	if err != nil {
		return err
	}

	return nil
}

func StartServer() (*http.Server, error) {
	err := config.LoadConfig()
	if err != nil {
		log.Println("could not load one or more config")
		return nil, err
	}
	//rLogger
	verbose := false
	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}
	rLogger := rLog.New(verbose)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	cfg := config.GetConfig()

	gracefulTimeout := time.Second * time.Duration(cfg.GracefulTimeout)

	db, err = infraMongo.New(ctx, cfg.DSN, cfg.Database, gracefulTimeout, infraMongo.SetLogger(rLogger))
	if err != nil {
		return nil, err
	}

	cache = infraCache.NewCacheDB(cfg.CacheURL, "")

	svc := service.SetupServiceConfig(cfg, db, cache, rLogger)

	log.Println("db initialized")

	server, err := api.Start(cfg, svc, rLogger)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}

	return server, nil
}

func StopServer(server *http.Server) error {
	defer db.Close(context.Background())
	defer cache.Client.Close()
	var err error
	graceful := func() error {
		log.Println("Shutting down server gracefully")
		return nil
	}

	forced := func() error {
		log.Println("Shutting down server forcefully")
		return nil
	}

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	errCh := make(chan error)
	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()
	if err = <-errCh; err != nil {
		log.Println(err)
		return err
	}

	err = api.Stop(server)
	if err != nil {
		log.Println("server stop err:", err)
		return err
	}

	return nil
}

// HandleSignals listen on the registered signals and fires the gracefulHandler for the
// first signal and the forceHandler (if any) for the next this function blocks and
// return any error that returned by any of the api first
func HandleSignals(sigs []os.Signal, gracefulHandler, forceHandler func() error) error {
	sigCh := make(chan os.Signal)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	grace := true
	for {
		select {
		case err := <-errCh:
			return err
		case <-sigCh:
			if grace {
				grace = false
				go func() {
					errCh <- gracefulHandler()
				}()
			} else if forceHandler != nil {
				err := forceHandler()
				errCh <- err
			}
		}
	}
}
