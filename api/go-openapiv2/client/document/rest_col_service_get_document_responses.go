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

// RestColServiceGetDocumentReader is a Reader for the RestColServiceGetDocument structure.
type RestColServiceGetDocumentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RestColServiceGetDocumentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRestColServiceGetDocumentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRestColServiceGetDocumentDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRestColServiceGetDocumentOK creates a RestColServiceGetDocumentOK with default headers values
func NewRestColServiceGetDocumentOK() *RestColServiceGetDocumentOK {
	return &RestColServiceGetDocumentOK{}
}

/*
RestColServiceGetDocumentOK describes a response with status code 200, with default header values.

A successful response.
*/
type RestColServiceGetDocumentOK struct {
	Payload *models.APIGetDocumentResponse
}

// IsSuccess returns true when this rest col service get document o k response has a 2xx status code
func (o *RestColServiceGetDocumentOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this rest col service get document o k response has a 3xx status code
func (o *RestColServiceGetDocumentOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this rest col service get document o k response has a 4xx status code
func (o *RestColServiceGetDocumentOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this rest col service get document o k response has a 5xx status code
func (o *RestColServiceGetDocumentOK) IsServerError() bool {
	return false
}

// IsCode returns true when this rest col service get document o k response a status code equal to that given
func (o *RestColServiceGetDocumentOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the rest col service get document o k response
func (o *RestColServiceGetDocumentOK) Code() int {
	return 200
}

func (o *RestColServiceGetDocumentOK) Error() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/docs/{documentId}][%d] restColServiceGetDocumentOK  %+v", 200, o.Payload)
}

func (o *RestColServiceGetDocumentOK) String() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/docs/{documentId}][%d] restColServiceGetDocumentOK  %+v", 200, o.Payload)
}

func (o *RestColServiceGetDocumentOK) GetPayload() *models.APIGetDocumentResponse {
	return o.Payload
}

func (o *RestColServiceGetDocumentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIGetDocumentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRestColServiceGetDocumentDefault creates a RestColServiceGetDocumentDefault with default headers values
func NewRestColServiceGetDocumentDefault(code int) *RestColServiceGetDocumentDefault {
	return &RestColServiceGetDocumentDefault{
		_statusCode: code,
	}
}

/*
RestColServiceGetDocumentDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type RestColServiceGetDocumentDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// IsSuccess returns true when this rest col service get document default response has a 2xx status code
func (o *RestColServiceGetDocumentDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this rest col service get document default response has a 3xx status code
func (o *RestColServiceGetDocumentDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this rest col service get document default response has a 4xx status code
func (o *RestColServiceGetDocumentDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this rest col service get document default response has a 5xx status code
func (o *RestColServiceGetDocumentDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this rest col service get document default response a status code equal to that given
func (o *RestColServiceGetDocumentDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the rest col service get document default response
func (o *RestColServiceGetDocumentDefault) Code() int {
	return o._statusCode
}

func (o *RestColServiceGetDocumentDefault) Error() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/docs/{documentId}][%d] RestColService_GetDocument default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceGetDocumentDefault) String() string {
	return fmt.Sprintf("[GET /v1/projects/{projectId}/collections/{collectionId}/docs/{documentId}][%d] RestColService_GetDocument default  %+v", o._statusCode, o.Payload)
}

func (o *RestColServiceGetDocumentDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *RestColServiceGetDocumentDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
