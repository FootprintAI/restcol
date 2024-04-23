package server

import (
	"context"

	grpcserver "github.com/sdinsure/agent/pkg/grpc/server"
	httpgateway "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"
	"github.com/sdinsure/agent/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/footprintai/restcol/api/pb/proto"
	appapp "github.com/footprintai/restcol/pkg/app"
)

type ServiceConfigure interface {
	apply(*ServiceConfig)
}

type ServiceConfig struct {
	unaryMiddlewares     []grpc.UnaryServerInterceptor
	streamMiddlewares    []grpc.StreamServerInterceptor
	transportCredentials credentials.TransportCredentials
}

func newServiceConfig(scs ...ServiceConfigure) *ServiceConfig {
	c := &ServiceConfig{}
	for _, sc := range scs {
		sc.apply(c)
	}
	return c
}

type middlewareConfigure struct {
	unaryMiddlewares  []grpc.UnaryServerInterceptor
	streamMiddlewares []grpc.StreamServerInterceptor
}

func (m middlewareConfigure) apply(sc *ServiceConfig) {
	sc.unaryMiddlewares = append(sc.unaryMiddlewares, m.unaryMiddlewares...)
	sc.streamMiddlewares = append(sc.streamMiddlewares, m.streamMiddlewares...)
}

func WithMiddlewareConfigure(unaryMiddlewares []grpc.UnaryServerInterceptor, streamMiddlewares []grpc.StreamServerInterceptor) middlewareConfigure {
	return middlewareConfigure{
		unaryMiddlewares:  unaryMiddlewares,
		streamMiddlewares: streamMiddlewares,
	}
}

type transportCredentialConfigure struct {
	transportCredentials credentials.TransportCredentials
}

func (t transportCredentialConfigure) apply(sc *ServiceConfig) {
	sc.transportCredentials = t.transportCredentials
}

func WithTransportCredential(tc credentials.TransportCredentials) transportCredentialConfigure {
	return transportCredentialConfigure{
		transportCredentials: tc,
	}
}

func NewServerService(
	grpcPort int,
	httpPort int,
	log logger.Logger,
	configs ...ServiceConfigure,
) (*ServerService, error) {

	config := newServiceConfig(configs...)

	svr := grpcserver.NewGrpcServerWithInterceptors(log, nil, config.unaryMiddlewares, config.streamMiddlewares, grpcserver.WithGrpcPort(grpcPort))
	httpGateway, err := httpgateway.NewHTTPGatewayServer(svr, log, httpPort, httpgateway.WithTransportCredentials(config.transportCredentials))
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
