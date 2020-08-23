package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/KentaKudo/go-todo-graphql/graph"
	"github.com/KentaKudo/go-todo-graphql/graph/generated"
	"github.com/KentaKudo/go-todo-graphql/internal/pb/service"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var gitHash = "overridden at compile time"

const (
	appName = "go-todo-graphql"
	appDesc = "The GraphQL wrapper to todo service"
)

func main() {
	app := cli.App(appName, appDesc)

	appPort := app.Int(cli.IntOpt{
		Desc:   "application http port",
		Name:   "app-port",
		EnvVar: "APP_PORT",
		Value:  8080,
	})

	todoAPI := app.String(cli.StringOpt{
		Name:   "todo-api",
		Desc:   "The Todo gRPC endpoint",
		EnvVar: "TODO_API",
		Value:  "localhost:8090",
	})

	logger := log.WithField("git_hash", gitHash)

	app.Action = func() {
		logger.Println("app start")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		todoAPIConn, err := initialiseGRPCClientConnection(ctx, *todoAPI)
		if err != nil {
			logger.WithError(err).Fatalln("initialise Todo gRPC client conn")
		}
		defer todoAPIConn.Close()

		todoAPICli := service.NewTodoAPIClient(todoAPIConn)
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{TodoClient: todoAPICli}}))

		router := mux.NewRouter()
		router.Handle("/query", srv)

		s := &http.Server{
			Handler: router,
			Addr:    net.JoinHostPort("", strconv.Itoa(*appPort)),
		}

		sigCh, errCh := make(chan os.Signal, 1), make(chan error, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := s.ListenAndServe(); err != nil {
				errCh <- fmt.Errorf("app server: %w", err)
			}
		}()

		select {
		case err := <-errCh:
			logger.WithError(err).Println("error received. attempt graceful shutdown.")
		case <-sigCh:
			logger.Println("termination signal received. attempt graceful shutdown.")
		}

		if err := s.Shutdown(context.Background()); err != nil {
			logger.WithError(err).Error("shutting down the http server")
		}

		cancel()
		wg.Wait()

		logger.Println("bye:)")
	}

	if err := app.Run(os.Args); err != nil {
		logger.WithError(err).Fatalln("app run")
	}
}

func initialiseGRPCClientConnection(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx,
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial, addr – %s: %w", addr, err)
	}

	return conn, nil
}
