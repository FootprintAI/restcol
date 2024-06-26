// Code generated by go-swagger; DO NOT EDIT.

package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/footprintai/restcol/api/go-openapiv2/models"
)

// RestColServiceGetSwaggerDoc2Reader is a Reader for the RestColServiceGetSwaggerDoc2 structure.
type RestColServiceGetSwaggerDoc2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceGetSwaggerDoc2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceGetSwaggerDoc2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceGetSwaggerDoc2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceGetSwaggerDoc2OK creates a RestColServiceGetSwaggerDoc2OK with default headers values
func NewRestColServiceGetSwaggerDoc2OK() *RestColServiceGetSwaggerDoc2OK {
	return &RestColServiceGetSwaggerDoc2OK{}
}

/*
RestColServiceGetSwaggerDoc2OK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceGetSwaggerDoc2OK struct {
	Payload *models.APIHTTPBody
}

// IsSuccess returns true when this rest col service get swagger doc2 o k response has a 2xx status code
func (o *RestColServiceGetSwaggerDoc2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service get swagger doc2 o k response has a 3xx status code
func (o *RestColServiceGetSwaggerDoc2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service get swagger doc2 o k response has a 4xx status code
func (o *RestColServiceGetSwaggerDoc2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service get swagger doc2 o k response has a 5xx status code
func (o *RestColServiceGetSwaggerDoc2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service get swagger doc2 o k response a status code equal to that given
func (o *RestColServiceGetSwaggerDoc2OK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service get swagger doc2 o k response
func (o *RestColServiceGetSwaggerDoc2OK) Code() int {
	return 200
}

func (o *RestColServiceGetSwaggerDoc2OK) Error() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/apidoc][%d] restColServiceGetSwaggerDoc2OK  %+v", 200, o.Payload)
}

func (o *RestColServiceGetSwaggerDoc2OK) String() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/apidoc][%d] restColServiceGetSwaggerDoc2OK  %+v", 200, o.Payload)
}

func (o *RestColServiceGetSwaggerDoc2OK) GetPayload() *models.APIHTTPBody {
	return o.Payload
}

func (o *RestColServiceGetSwaggerDoc2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIHTTPBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceGetSwaggerDoc2Default creates a RestColServiceGetSwaggerDoc2Default with default headers values
func NewRestColServiceGetSwaggerDoc2Default(code int) *RestColServiceGetSwaggerDoc2Default {
	return &RestColServiceGetSwaggerDoc2Default{
		_statusCode: code,
	}
}

/*
RestColServiceGetSwaggerDoc2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceGetSwaggerDoc2Default struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service get swagger doc2 default response has a 2xx status code
func (o *RestColServiceGetSwaggerDoc2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service get swagger doc2 default response has a 3xx status code
func (o *RestColServiceGetSwaggerDoc2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service get swagger doc2 default response has a 4xx status code
func (o *RestColServiceGetSwaggerDoc2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service get swagger doc2 default response has a 5xx status code
func (o *RestColServiceGetSwaggerDoc2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service get swagger doc2 default response a status code equal to that given
func (o *RestColServiceGetSwaggerDoc2Default) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service get swagger doc2 default response
func (o *RestColServiceGetSwaggerDoc2Default) Code() int {
	return o._statusCode
}

func (o *RestColServiceGetSwaggerDoc2Default) Error() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/apidoc][%d] RestColService_GetSwaggerDoc2 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceGetSwaggerDoc2Default) String() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/apidoc][%d] RestColService_GetSwaggerDoc2 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceGetSwaggerDoc2Default) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceGetSwaggerDoc2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
