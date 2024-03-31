package main

import (
	"context"
	"flag"

	logger "github.com/sdinsure/agent/pkg/logger"
	storageflags "github.com/sdinsure/agent/pkg/storage/flags"
	postgresstorage "github.com/sdinsure/agent/pkg/storage/postgres"

	appapp "github.com/footprintai/restcol/pkg/app"
	appauthn "github.com/footprintai/restcol/pkg/authn"
	appauthz "github.com/footprintai/restcol/pkg/authz"
	dummy "github.com/footprintai/restcol/pkg/dummy"
	appserver "github.com/footprintai/restcol/pkg/server"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
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
	projectCURD := projectsstorage.NewProjectCURD(postgresDb)
	if err := projectCURD.AutoMigrate(); err != nil {
		log.Fatal("restcol: collection automigrate failed, err:%+v\n", err)
	}
	collectionCURD := collectionsstorage.NewCollectionCURD(postgresDb)
	if err := collectionCURD.AutoMigrate(); err != nil {
		log.Fatal("restcol: collection automigrate failed, err:%+v\n", err)
	}

	documentCURD := documentsstorage.NewDocumentCURD(postgresDb)
	if err := documentCURD.AutoMigrate(); err != nil {
		log.Fatal("restcol: document automigrate failed, err:%+v\n", err)
	}
	dummyProject := dummy.NewDummyProject(projectCURD)
	if err = dummyProject.Init(context.Background()); err != nil {
		log.Fatal("failed to init dummy project, err:%s\n", err)
	}

	authZKeeper := &appauthz.AllowEveryOne{}
	authNParser := &appauthn.AnnonymousClaimParser{}
	svr, err := appserver.NewServerService(*grpcPort, *httpPort, log, authZKeeper, authNParser)
	if err != nil {
		log.Fatal("failed to start server, err:%+v\n", err)
	}
	svr.AddGatewayRoutes()
	app := appapp.NewRestColServiceServerService(
		log,
		collectionCURD,
		documentCURD,
		dummyProject,
	)
	appserver.RegisterService(svr, app)
	svr.Start()
}
