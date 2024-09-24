package configuration

import (
	"fmt"
	"time"
)

type GRPC struct {
	Port    int
	Timeout time.Duration
}

type SQLDataBase struct {
	DriverName string
	Username   string
	Password   string
	Host       string
	Port       int
	DBName     string
}

type Configuration struct {
	Env      Env
	Host     string
	Port     int
	GRPC     GRPC
	DataBase SQLDataBase
}

func (grpc *GRPC) String() string {
	return fmt.Sprintf("GRPC: {Port: %v, Timeout: %v}", grpc.Port, grpc.Timeout)
}

func (db *SQLDataBase) String() string {
	return fmt.Sprintf("DataBase: %v://%v:***@%v:%v/%v",
		db.DriverName,
		db.Username,
		db.Host,
		db.Port,
		db.DBName,
	)
}

func (cfg Configuration) String() string {
	return fmt.Sprintf(
		"---\n"+
			"Env: %v\n"+
			"Server: %v:%v\n"+
			"%v\n"+
			"%v\n"+
			"---\n",
		cfg.Env,
		cfg.Host,
		cfg.Port,
		cfg.GRPC,
		cfg.DataBase,
	)
}
