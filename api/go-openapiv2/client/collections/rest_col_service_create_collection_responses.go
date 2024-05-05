// Code generated by go-swagger; DO NOT EDIT.

package collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/footprintai/restcol/api/go-openapiv2/models"
)

// RestColServiceCreateCollectionReader is a Reader for the RestColServiceCreateCollection structure.
type RestColServiceCreateCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceCreateCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceCreateCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceCreateCollectionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceCreateCollectionOK creates a RestColServiceCreateCollectionOK with default headers values
func NewRestColServiceCreateCollectionOK() *RestColServiceCreateCollectionOK {
	return &RestColServiceCreateCollectionOK{}
}

/*
RestColServiceCreateCollectionOK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceCreateCollectionOK struct {
	Payload *models.APICreateCollectionResponse
}

// IsSuccess returns true when this rest col service create collection o k response has a 2xx status code
func (o *RestColServiceCreateCollectionOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service create collection o k response has a 3xx status code
func (o *RestColServiceCreateCollectionOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service create collection o k response has a 4xx status code
func (o *RestColServiceCreateCollectionOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service create collection o k response has a 5xx status code
func (o *RestColServiceCreateCollectionOK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service create collection o k response a status code equal to that given
func (o *RestColServiceCreateCollectionOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service create collection o k response
func (o *RestColServiceCreateCollectionOK) Code() int {
	return 200
}

func (o *RestColServiceCreateCollectionOK) Error() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections][%d] restColServiceCreateCollectionOK  %+v", 200, o.Payload)
}

func (o *RestColServiceCreateCollectionOK) String() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections][%d] restColServiceCreateCollectionOK  %+v", 200, o.Payload)
}

func (o *RestColServiceCreateCollectionOK) GetPayload() *models.APICreateCollectionResponse {
	return o.Payload
}

func (o *RestColServiceCreateCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APICreateCollectionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceCreateCollectionDefault creates a RestColServiceCreateCollectionDefault with default headers values
func NewRestColServiceCreateCollectionDefault(code int) *RestColServiceCreateCollectionDefault {
	return &RestColServiceCreateCollectionDefault{
		_statusCode: code,
	}
}

/*
RestColServiceCreateCollectionDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceCreateCollectionDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service create collection default response has a 2xx status code
func (o *RestColServiceCreateCollectionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service create collection default response has a 3xx status code
func (o *RestColServiceCreateCollectionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service create collection default response has a 4xx status code
func (o *RestColServiceCreateCollectionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service create collection default response has a 5xx status code
func (o *RestColServiceCreateCollectionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service create collection default response a status code equal to that given
func (o *RestColServiceCreateCollectionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service create collection default response
func (o *RestColServiceCreateCollectionDefault) Code() int {
	return o._statusCode
}

func (o *RestColServiceCreateCollectionDefault) Error() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections][%d] RestColService_CreateCollection default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceCreateCollectionDefault) String() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections][%d] RestColService_CreateCollection default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceCreateCollectionDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceCreateCollectionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
