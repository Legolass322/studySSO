package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"sso/internal/db"
	authgrpc "sso/internal/grpc/auth"

	serviceAuth "sso/internal/services/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
	storage *db.Database,
) *App {
	gRPCServer := grpc.NewServer()

	authService := serviceAuth.New(log, storage, storage, storage)

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Runs server and panics if any error occurs
func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(fmt.Errorf("Couldn't run: %w", err))
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", int(app.port)),
	)

	log.Info("Listen grpc")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("Error trying to listen [%s]: %w", op, err)
	}

	log.Info("gRPC is running", slog.String("addr", listener.Addr().String()))

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("Error while listening [%s]: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpc.Stop"

	app.log.With(slog.String("op", op)).Info("Stopping grpc...")

	app.gRPCServer.GracefulStop()
}
