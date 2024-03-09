package main

import (
	"flag"

	logger "github.com/sdinsure/agent/pkg/logger"
	storageflags "github.com/sdinsure/agent/pkg/storage/flags"
	postgresstorage "github.com/sdinsure/agent/pkg/storage/postgres"

	appapp "github.com/footprintai/restcol/pkg/app"
	appauthn "github.com/footprintai/restcol/pkg/authn"
	appauthz "github.com/footprintai/restcol/pkg/authz"
	appserver "github.com/footprintai/restcol/pkg/server"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
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
		log.Fatal("failed to init dsn, err:%s\n", err)
	}
	_ = postgresDb
	collectionCURD := collectionsstorage.NewCollectionCURD(postgresDb)
	if err := collectionCURD.AutoMigrate(); err != nil {
		log.Fatal("restcol: collection automigrate failed, err:%+v\n", err)
	}
	authZKeeper := &appauthz.AllowEveryOne{}
	authNParser := &appauthn.AnnonymousClaimParser{}
	svr, err := appserver.NewServerService(*grpcPort, *httpPort, log, authZKeeper, authNParser)
	if err != nil {
		log.Fatal("failed to start server, err:%+v\n", err)
	}
	svr.AddGatewayRoutes()
	app := appapp.NewRestColServiceServerService(log, collectionCURD)
	appserver.RegisterService(svr, app)
	svr.Start()
}
