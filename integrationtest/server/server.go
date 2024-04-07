package integrationtestserver

import (
	"context"

	"github.com/sdinsure/agent/pkg/logger"

	appapp "github.com/footprintai/restcol/pkg/app"
	appauthn "github.com/footprintai/restcol/pkg/authn"
	appauthz "github.com/footprintai/restcol/pkg/authz"
	dummy "github.com/footprintai/restcol/pkg/dummy"
	appserver "github.com/footprintai/restcol/pkg/server"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
	postgresstorage "github.com/sdinsure/agent/pkg/storage/postgres"
)

type StartStopServer interface {
	Start() error
	Stop() error
}

var (
	_ StartStopServer = &Server{}
)

func NewServer(
	grpcPort int,
	httpPort int,
	postgresDb *postgresstorage.PostgresDb,
	log logger.Logger,
) (*Server, error) {

	svr, err := makeServerService(grpcPort, httpPort, postgresDb, log)
	if err != nil {
		return nil, err
	}
	return &Server{
		server:   svr,
		httpPort: httpPort,
	}, nil
}

func makeServerService(
	grpcPort int,
	httpPort int,
	postgresDb *postgresstorage.PostgresDb,
	log logger.Logger,
) (*appserver.ServerService, error) {
	projectCURD := projectsstorage.NewProjectCURD(postgresDb)
	if err := projectCURD.AutoMigrate(); err != nil {
		return nil, err
	}
	collectionCURD := collectionsstorage.NewCollectionCURD(postgresDb)
	if err := collectionCURD.AutoMigrate(); err != nil {
		return nil, err
	}

	documentCURD := documentsstorage.NewDocumentCURD(postgresDb)
	if err := documentCURD.AutoMigrate(); err != nil {
		return nil, err
	}
	dummyProject := dummy.NewDummyProject(projectCURD)
	if err := dummyProject.Init(context.Background()); err != nil {
		return nil, err
	}

	authZKeeper := &appauthz.AllowEveryOne{}
	authNParser := &appauthn.AnnonymousClaimParser{}
	svr, err := appserver.NewServerService(grpcPort, httpPort, log, authZKeeper, authNParser)
	if err != nil {
		return nil, err
	}
	svr.AddGatewayRoutes()
	app := appapp.NewRestColServiceServerService(
		log,
		collectionCURD,
		documentCURD,
		dummyProject,
	)
	appserver.RegisterService(svr, app)
	return svr, nil
}

type Server struct {
	server   *appserver.ServerService
	httpPort int
}

// Start starts Server and blocks forever
func (s *Server) Start() error {
	return s.server.Start()
}

// Stop stops server
func (s *Server) Stop() error {
	return s.server.Stop()
}
