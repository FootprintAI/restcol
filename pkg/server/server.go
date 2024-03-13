package server

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcserver "github.com/sdinsure/agent/pkg/grpc/server"
	httpgateway "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"
	authnmiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/authn"
	authzmiddleware "github.com/sdinsure/agent/pkg/grpc/server/middleware/authz"
	"github.com/sdinsure/agent/pkg/logger"
	"google.golang.org/grpc"

	pb "github.com/footprintai/restcol/api/pb/proto"
	appapp "github.com/footprintai/restcol/pkg/app"
	appauthz "github.com/footprintai/restcol/pkg/authz"
	appmiddleware "github.com/footprintai/restcol/pkg/middleware"
)

func NewServerService(
	grpcPort int,
	httpPort int,
	log logger.Logger,
	authzService appauthz.AuthzService,
	authnClaimParser authnmiddleware.ClaimParser,
) (*ServerService, error) {

	authNMiddleware := authnmiddleware.NewAuthNMiddleware(
		log,
		authnClaimParser,
	)

	authZMiddleware := authzmiddleware.NewAuthZMiddleware(
		log,
		appmiddleware.NewAuthzMiddlwareAdaptor(authzService),
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
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_auth.UnaryServerInterceptor(authNMiddleware.AuthFunc),
		grpc_auth.UnaryServerInterceptor(authZMiddleware.AuthFunc),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_auth.StreamServerInterceptor(authNMiddleware.AuthFunc),
		grpc_auth.StreamServerInterceptor(authZMiddleware.AuthFunc),
	}

	svr := grpcserver.NewGrpcServerWithInterceptors(grpcPort, log, nil, unaryInterceptors, streamInterceptors)
	httpGateway, err := httpgateway.NewHTTPGatewayServer(svr, log, httpPort)
	if err != nil {
		return nil, err
	}
	return &ServerService{
		log:         log,
		svr:         svr,
		httpGateway: httpGateway,
	}, nil
}

type ServerService struct {
	log         logger.Logger
	svr         *grpcserver.GrpcServer
	httpGateway *httpgateway.HTTPGatewayServer
}

func (s *ServerService) Start() error {
	go func() {
		s.httpGateway.ListenAndServe()
	}()

	return s.httpGateway.WaitForSIGTERM()
}

func (s *ServerService) Stop() error {
	return s.httpGateway.Shutdown(context.Background())
}

func RegisterService(svr *ServerService, app *appapp.RestColServiceServerService) {
	pb.RegisterRestColServiceServer(svr.svr, app)
	svr.httpGateway.RegisterHandlers(pb.RegisterRestColServiceHandler)
}
