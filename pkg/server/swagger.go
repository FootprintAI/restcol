package server

import (
	grpchttpgatewayserver "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"

	api "github.com/footprintai/restcol/api"
)

func (s *ServerService) AddGatewayRoutes() error {
	return s.httpGateway.AddRoutes(
		grpchttpgatewayserver.NewSwaggerRoute(),
		grpchttpgatewayserver.NewOpenAPIV2Route("/openapiv2/", api.OpenApiV2HttpHandler),
	)
}
