package integrationtestclient

import (
	"io"

	restcolopenapi "github.com/footprintai/restcol/api/go-openapiv2/client"
	restcolopenapiswagger "github.com/footprintai/restcol/api/go-openapiv2/client/swagger"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
	"github.com/go-openapi/runtime"
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

type RawSwaggerDocReader struct {
	restColServiceGetSwaggerDocReader restcolopenapiswagger.RestColServiceGetSwaggerDocReader
}

func (r *RawSwaggerDocReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := restcolopenapiswagger.NewRestColServiceGetSwaggerDocOK()
		result.Payload = new(restcolopenapimodel.APIHTTPBody)
		result.Payload.Data, _ = io.ReadAll(response.Body())
		result.Payload.ContentType = response.GetHeader("content-type")
		return result, nil
	default:
		return r.restColServiceGetSwaggerDocReader.ReadResponse(response, consumer)
	}
}

var (
	_ runtime.ClientResponseReader = &RawSwaggerDocReader{}
)

func WithRawSwaggerDocReader() restcolopenapiswagger.ClientOption {
	return restcolopenapiswagger.ClientOption(func(co *runtime.ClientOperation) {
		co.Reader = &RawSwaggerDocReader{}
	})
}
