package server

import (
	"embed"
	"net/http"
	"strings"

	grpchttpgatewayserver "github.com/sdinsure/agent/pkg/grpc/server/httpgateway"

	api "github.com/footprintai/restcol/api"
)

type GatewayRouteAdder interface {
	AddGatewayRoutes(routes ...*grpchttpgatewayserver.Route) error
}

// Only the initializer. The rest of swagger-ui-dist - the bundles, the CSS, the
// fonts, several megabytes of it - stays in the agent's embed, because we have
// no reason to carry a second copy and every reason not to let the two drift.
//
//go:embed swaggerui/swagger-initializer.js
var initializerAssets embed.FS

const initializerPath = "swagger-initializer.js"

// swaggerUIHandler serves the agent's swagger-ui assets, with one file
// replaced.
//
// The agent embeds swagger-ui-dist verbatim, and dist's swagger-initializer.js
// points at https://petstore.swagger.io/v2/swagger.json. Its own comment says
// the line "will be replaced by docker/configurator, when it runs in a
// docker-container" - we embed the assets rather than running that container,
// so nothing replaced it, and every deployment has been serving Swagger's demo
// API where ours should be.
//
// Overriding one file rather than vendoring the whole directory keeps the fix
// the size of the defect.
func swaggerUIHandler() http.Handler {
	assets := http.FileServer(http.FS(initializerAssets))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/"+initializerPath) ||
			strings.TrimPrefix(r.URL.Path, "/") == initializerPath {
			// http.FS serves from the embed root, where the file lives under
			// swaggerui/. Ask for it by that name rather than by whatever
			// prefix the request arrived on.
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/swaggerui/" + initializerPath
			assets.ServeHTTP(w, r2)
			return
		}
		grpchttpgatewayserver.NewSwaggerRoute().Handler.ServeHTTP(w, r)
	})
}

func AddSwaggerRoutes(s GatewayRouteAdder) error {
	return s.AddGatewayRoutes(
		&grpchttpgatewayserver.Route{
			Pattern: "/swaggerui/",
			Handler: swaggerUIHandler(),
		},
		grpchttpgatewayserver.NewOpenAPIV2Route("/openapiv2/", api.OpenApiV2HttpHandler),
	)
}
