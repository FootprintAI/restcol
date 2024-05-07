package integrationtestclient

import (
	"io"

	"github.com/go-openapi/runtime"

	restcolopenapiswagger "github.com/footprintai/restcol/api/go-openapiv2/client/swagger"
	restcolopenapimodel "github.com/footprintai/restcol/api/go-openapiv2/models"
)

type RawSwaggerDocReader struct {
	restColServiceGetSwaggerDocReader restcolopenapiswagger.RestColServiceGetSwaggerDocReader
}

func (r *RawSwaggerDocReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := restcolopenapiswagger.NewRestColServiceGetSwaggerDoc2OK()
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
