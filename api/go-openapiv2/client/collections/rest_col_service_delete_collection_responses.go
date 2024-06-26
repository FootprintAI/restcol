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

// RestColServiceDeleteCollectionReader is a Reader for the RestColServiceDeleteCollection structure.
type RestColServiceDeleteCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceDeleteCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceDeleteCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceDeleteCollectionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceDeleteCollectionOK creates a RestColServiceDeleteCollectionOK with default headers values
func NewRestColServiceDeleteCollectionOK() *RestColServiceDeleteCollectionOK {
	return &RestColServiceDeleteCollectionOK{}
}

/*
RestColServiceDeleteCollectionOK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceDeleteCollectionOK struct {
	Payload models.APIDeleteCollectionResponse
}

// IsSuccess returns true when this rest col service delete collection o k response has a 2xx status code
func (o *RestColServiceDeleteCollectionOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service delete collection o k response has a 3xx status code
func (o *RestColServiceDeleteCollectionOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service delete collection o k response has a 4xx status code
func (o *RestColServiceDeleteCollectionOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service delete collection o k response has a 5xx status code
func (o *RestColServiceDeleteCollectionOK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service delete collection o k response a status code equal to that given
func (o *RestColServiceDeleteCollectionOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service delete collection o k response
func (o *RestColServiceDeleteCollectionOK) Code() int {
	return 200
}

func (o *RestColServiceDeleteCollectionOK) Error() string {
	return fmt.Sprintf("[DELETE /v1/projects/{projectId}/collections/{collectionId}][%d] restColServiceDeleteCollectionOK  %+v", 200, o.Payload)
}

func (o *RestColServiceDeleteCollectionOK) String() string {
	return fmt.Sprintf("[DELETE /v1/projects/{projectId}/collections/{collectionId}][%d] restColServiceDeleteCollectionOK  %+v", 200, o.Payload)
}

func (o *RestColServiceDeleteCollectionOK) GetPayload() models.APIDeleteCollectionResponse {
	return o.Payload
}

func (o *RestColServiceDeleteCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceDeleteCollectionDefault creates a RestColServiceDeleteCollectionDefault with default headers values
func NewRestColServiceDeleteCollectionDefault(code int) *RestColServiceDeleteCollectionDefault {
	return &RestColServiceDeleteCollectionDefault{
		_statusCode: code,
	}
}

/*
RestColServiceDeleteCollectionDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceDeleteCollectionDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service delete collection default response has a 2xx status code
func (o *RestColServiceDeleteCollectionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service delete collection default response has a 3xx status code
func (o *RestColServiceDeleteCollectionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service delete collection default response has a 4xx status code
func (o *RestColServiceDeleteCollectionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service delete collection default response has a 5xx status code
func (o *RestColServiceDeleteCollectionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service delete collection default response a status code equal to that given
func (o *RestColServiceDeleteCollectionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service delete collection default response
func (o *RestColServiceDeleteCollectionDefault) Code() int {
	return o._statusCode
}

func (o *RestColServiceDeleteCollectionDefault) Error() string {
	return fmt.Sprintf("[DELETE /v1/projects/{projectId}/collections/{collectionId}][%d] RestColService_DeleteCollection default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceDeleteCollectionDefault) String() string {
	return fmt.Sprintf("[DELETE /v1/projects/{projectId}/collections/{collectionId}][%d] RestColService_DeleteCollection default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceDeleteCollectionDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceDeleteCollectionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
