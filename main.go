package main

import (
	"flag"
	"time"

	logger "github.com/sdinsure/agent/pkg/logger"
	storageflags "github.com/sdinsure/agent/pkg/storage/flags"
	postgresstorage "github.com/sdinsure/agent/pkg/storage/postgres"

	integrationtestserver "github.com/footprintai/restcol/integrationtest/server"
	"github.com/footprintai/restcol/pkg/version"
)

var grpcPort = flag.Int("grpc_port", 50090, "The server grpc port")
var httpPort = flag.Int("http_port", 50091, "The server http port")
var restcolPostgresFlags = storageflags.NewPrefixFlagSet("restcol")

func main() {

	restcolPostgresFlags.Init()

	flag.Parse()
	version.Print()

	log := logger.NewLogger()

	postgresDb, err := postgresstorage.NewPostgresDbHelper(log, restcolPostgresFlags)
	if err != nil {
		log.Fatal("failed to init dsn, err:%+v\n", err)
	}

	svr, err := integrationtestserver.NewServer(*grpcPort, *httpPort, postgresDb, log)
	if err != nil {
		log.Fatal("failed to init dsn, err:%+v\n", err)
	}
	if err := svr.Start(30 * time.Second); err != nil {
		log.Fatal("failed to launch server:%+v\n", err)
	}
}
