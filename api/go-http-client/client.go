package gohttpclient

import (
	restcolopenapi "github.com/footprintai/restcol/api/go-openapiv2/client"
	"github.com/sdinsure/agent/pkg/http/openapi"
)

func MustNewClient(endpoint string) *restcolopenapi.RestColAPIDocumentations {
	return restcolopenapi.New(
		openapi.MustNew(
			endpoint,
			restcolopenapi.DefaultBasePath,
		),
		nil,
	)
}
