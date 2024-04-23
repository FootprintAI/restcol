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

// RestColServiceQueryDocument2Reader is a Reader for the RestColServiceQueryDocument2 structure.
type RestColServiceQueryDocument2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceQueryDocument2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceQueryDocument2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceQueryDocument2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceQueryDocument2OK creates a RestColServiceQueryDocument2OK with default headers values
func NewRestColServiceQueryDocument2OK() *RestColServiceQueryDocument2OK {
	return &RestColServiceQueryDocument2OK{}
}

/*
RestColServiceQueryDocument2OK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceQueryDocument2OK struct {
	Payload *models.APIQueryDocumentResponse
}

// IsSuccess returns true when this rest col service query document2 o k response has a 2xx status code
func (o *RestColServiceQueryDocument2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service query document2 o k response has a 3xx status code
func (o *RestColServiceQueryDocument2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service query document2 o k response has a 4xx status code
func (o *RestColServiceQueryDocument2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service query document2 o k response has a 5xx status code
func (o *RestColServiceQueryDocument2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service query document2 o k response a status code equal to that given
func (o *RestColServiceQueryDocument2OK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service query document2 o k response
func (o *RestColServiceQueryDocument2OK) Code() int {
	return 200
}

func (o *RestColServiceQueryDocument2OK) Error() string {
	return fmt.Sprintf("[GET /v1/collections/{collectionId}:query][%d] restColServiceQueryDocument2OK  %+v", 200, o.Payload)
}

func (o *RestColServiceQueryDocument2OK) String() string {
	return fmt.Sprintf("[GET /v1/collections/{collectionId}:query][%d] restColServiceQueryDocument2OK  %+v", 200, o.Payload)
}

func (o *RestColServiceQueryDocument2OK) GetPayload() *models.APIQueryDocumentResponse {
	return o.Payload
}

func (o *RestColServiceQueryDocument2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIQueryDocumentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceQueryDocument2Default creates a RestColServiceQueryDocument2Default with default headers values
func NewRestColServiceQueryDocument2Default(code int) *RestColServiceQueryDocument2Default {
	return &RestColServiceQueryDocument2Default{
		_statusCode: code,
	}
}

/*
RestColServiceQueryDocument2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceQueryDocument2Default struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service query document2 default response has a 2xx status code
func (o *RestColServiceQueryDocument2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service query document2 default response has a 3xx status code
func (o *RestColServiceQueryDocument2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service query document2 default response has a 4xx status code
func (o *RestColServiceQueryDocument2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service query document2 default response has a 5xx status code
func (o *RestColServiceQueryDocument2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service query document2 default response a status code equal to that given
func (o *RestColServiceQueryDocument2Default) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service query document2 default response
func (o *RestColServiceQueryDocument2Default) Code() int {
	return o._statusCode
}

func (o *RestColServiceQueryDocument2Default) Error() string {
	return fmt.Sprintf("[GET /v1/collections/{collectionId}:query][%d] RestColService_QueryDocument2 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceQueryDocument2Default) String() string {
	return fmt.Sprintf("[GET /v1/collections/{collectionId}:query][%d] RestColService_QueryDocument2 default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceQueryDocument2Default) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceQueryDocument2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}