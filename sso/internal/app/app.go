package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/db"
	cfg "sso/pkg/configuration"
)

type App struct {
	GRPCApp *grpcapp.App
	Storage *db.Database
}

func New(
	log *slog.Logger,
	cfg cfg.Configuration,
) *App {
	storage := db.MustNew(log, (*db.DBCfg)(&cfg.DataBase))

	grpcApplication := grpcapp.New(log, cfg.GRPC.Port, storage)

	return &App{
		GRPCApp: grpcApplication,
		Storage: storage,
	}
}