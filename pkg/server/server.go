package server

import (
	"context"

	grpcserver "github.com/sdinsure/agent/pkg/grpc/server"
	httpgateway "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"
	"github.com/sdinsure/agent/pkg/logger"
	"google.golang.org/grpc"

	pb "github.com/footprintai/restcol/api/pb/proto"
	appapp "github.com/footprintai/restcol/pkg/app"
)

func NewServerService(
	grpcPort int,
	httpPort int,
	log logger.Logger,
	unaryMiddlewares []grpc.UnaryServerInterceptor,
	streamMiddlewares []grpc.StreamServerInterceptor,
) (*ServerService, error) {
	svr := grpcserver.NewGrpcServerWithInterceptors(log, nil, unaryMiddlewares, streamMiddlewares, grpcserver.WithGrpcPort(grpcPort))
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
