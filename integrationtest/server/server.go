package integrationtestserver

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"

	appapp "github.com/footprintai/restcol/pkg/app"
	appauthn "github.com/footprintai/restcol/pkg/authn"
	appauthz "github.com/footprintai/restcol/pkg/authz"
	dummy "github.com/footprintai/restcol/pkg/dummy"
	appmiddleware "github.com/footprintai/restcol/pkg/middleware"
	runtimeprojectgetter "github.com/footprintai/restcol/pkg/runtime/getter"
	appserver "github.com/footprintai/restcol/pkg/server"
	collectionsstorage "github.com/footprintai/restcol/pkg/storage/collections"
	documentsstorage "github.com/footprintai/restcol/pkg/storage/documents"
	projectsstorage "github.com/footprintai/restcol/pkg/storage/projects"
	authnmiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/authn"
	authzmiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/authz"
	identitymiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/identity"
	"github.com/sdinsure/agent/pkg/logger"
	sdinsureruntime "github.com/sdinsure/agent/pkg/runtime"
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

	projectResolver := sdinsureruntime.NewProjectResolver(log, runtimeprojectgetter.NewRuntimeProjectGetter(projectCURD))

	authNMiddleware := authnmiddleware.NewAuthNMiddleware(
		log,
		&appauthn.AnnonymousClaimParser{},
		authnmiddleware.EnableAnnonymous(true),
	)
	//authnMiddleware := grpc_auth.UnaryServerInterceptor(.AuthFunc)
	authZMiddleware := authzmiddleware.NewAuthZMiddleware(
		log,
		appmiddleware.NewAuthzMiddlwareAdaptor(&appauthz.AllowEveryOne{}),
		//middleware.WithSkippedAuthZPaths([]middleware.HttpPath{
		//	middleware.HttpPath{
		//		RawPath:   "/v1/login",
		//		RawMethod: http.MethodPost,
		//	},
		//	middleware.HttpPath{
		//		RawPath:   "/v1/user",
		//		RawMethod: http.MethodPost,
		//	},
		//}),
	)
	projectIdentityMiddleware := identitymiddleware.NewProjectIdentityMiddleware(projectResolver)

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_auth.UnaryServerInterceptor(authNMiddleware.AuthFunc),
		grpc_auth.UnaryServerInterceptor(authZMiddleware.AuthFunc),
		projectIdentityMiddleware.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_auth.StreamServerInterceptor(authNMiddleware.AuthFunc),
		grpc_auth.StreamServerInterceptor(authZMiddleware.AuthFunc),
		projectIdentityMiddleware.StreamServerInterceptor(),
	}

	svr, err := appserver.NewServerService(grpcPort, httpPort, log, unaryInterceptors, streamInterceptors)
	if err != nil {
		return nil, err
	}
	svr.AddGatewayRoutes()
	app := appapp.NewRestColServiceServerService(
		log,
		collectionCURD,
		documentCURD,
	)
	app.SetDefaultProjectResolver(projectResolver)
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
