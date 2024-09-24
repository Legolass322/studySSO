package cfg

import (
	cfg "sso/pkg/configuration"
	"time"
)

var config *cfg.Configuration

func InitConfiguration(env cfg.Env) cfg.Configuration {
	var meta = cfg.GetMetaConfiguration()

	c := meta.SetSlice(cfg.InputSlice{Env: cfg.Production}, cfg.Configuration{
		Env:  cfg.Production,
		Host: "0.0.0.0",
		Port: 8080,
		GRPC: cfg.GRPC{
			Port: 44044,
			Timeout: time.Minute, // todo
		},
		DataBase: cfg.SQLDataBase{ // todo
			DriverName: "postgres",
			Username: "username",
			Password: "password",
			Host: "localhost",
			Port: 5432,
			DBName: "sso",
		},
	}).InheritSlice(
		cfg.InputSlice{Env: cfg.Testing},
		cfg.InputSlice{Env: cfg.Production},
		func(config cfg.Configuration) cfg.Configuration {
			config.Env = cfg.Testing
			return config
		},
	).InheritSlice(
		cfg.InputSlice{Env: cfg.Local},
		cfg.InputSlice{Env: cfg.Testing},
		func(config cfg.Configuration) cfg.Configuration {
			config.Env = cfg.Local
			return config
		},
	).GetConfiguration(cfg.InputSlice{Env: env})

	config = &c
	return c
}

func Get() *cfg.Configuration {
	if config == nil {
		panic("Config not initialized")
	}

	return config
}
