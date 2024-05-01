package server

import (
	grpchttpgatewayserver "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"

	api "github.com/footprintai/restcol/api"
)

type GatewayRouteAdder interface {
	AddGatewayRoutes(routes ...*grpchttpgatewayserver.Route) error
}

func AddSwaggerRoutes(s GatewayRouteAdder) error {
	return s.AddGatewayRoutes(
		grpchttpgatewayserver.NewSwaggerRoute(),
		grpchttpgatewayserver.NewOpenAPIV2Route("/openapiv2/", api.OpenApiV2HttpHandler),
	)
}
