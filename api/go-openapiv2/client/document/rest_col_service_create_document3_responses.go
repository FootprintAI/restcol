// Code generated by go-swagger; DO NOT EDIT.

package document

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/footprintai/restcol/api/go-openapiv2/models"
)

// RestColServiceCreateDocument3Reader is a Reader for the RestColServiceCreateDocument3 structure.
type RestColServiceCreateDocument3Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceCreateDocument3Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceCreateDocument3OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceCreateDocument3Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceCreateDocument3OK creates a RestColServiceCreateDocument3OK with default headers values
func NewRestColServiceCreateDocument3OK() *RestColServiceCreateDocument3OK {
	return &RestColServiceCreateDocument3OK{}
}

/*
RestColServiceCreateDocument3OK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceCreateDocument3OK struct {
	Payload *models.APICreateDocumentResponse
}

// IsSuccess returns true when this rest col service create document3 o k response has a 2xx status code
func (o *RestColServiceCreateDocument3OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service create document3 o k response has a 3xx status code
func (o *RestColServiceCreateDocument3OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service create document3 o k response has a 4xx status code
func (o *RestColServiceCreateDocument3OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service create document3 o k response has a 5xx status code
func (o *RestColServiceCreateDocument3OK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service create document3 o k response a status code equal to that given
func (o *RestColServiceCreateDocument3OK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service create document3 o k response
func (o *RestColServiceCreateDocument3OK) Code() int {
	return 200
}

func (o *RestColServiceCreateDocument3OK) Error() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections/{collectionId}:newdoc][%d] restColServiceCreateDocument3OK  %+v", 200, o.Payload)
}

func (o *RestColServiceCreateDocument3OK) String() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections/{collectionId}:newdoc][%d] restColServiceCreateDocument3OK  %+v", 200, o.Payload)
}

func (o *RestColServiceCreateDocument3OK) GetPayload() *models.APICreateDocumentResponse {
	return o.Payload
}

func (o *RestColServiceCreateDocument3OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APICreateDocumentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceCreateDocument3Default creates a RestColServiceCreateDocument3Default with default headers values
func NewRestColServiceCreateDocument3Default(code int) *RestColServiceCreateDocument3Default {
	return &RestColServiceCreateDocument3Default{
		_statusCode: code,
	}
}

/*
RestColServiceCreateDocument3Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceCreateDocument3Default struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service create document3 default response has a 2xx status code
func (o *RestColServiceCreateDocument3Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service create document3 default response has a 3xx status code
func (o *RestColServiceCreateDocument3Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service create document3 default response has a 4xx status code
func (o *RestColServiceCreateDocument3Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service create document3 default response has a 5xx status code
func (o *RestColServiceCreateDocument3Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service create document3 default response a status code equal to that given
func (o *RestColServiceCreateDocument3Default) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service create document3 default response
func (o *RestColServiceCreateDocument3Default) Code() int {
	return o._statusCode
}

func (o *RestColServiceCreateDocument3Default) Error() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections/{collectionId}:newdoc][%d] RestColService_CreateDocument3 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceCreateDocument3Default) String() string {
	return fmt.Sprintf("[POST /v1/projects/{projectId}/collections/{collectionId}:newdoc][%d] RestColService_CreateDocument3 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceCreateDocument3Default) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceCreateDocument3Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
