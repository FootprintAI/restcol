// Code generated by go-swagger; DO NOT EDIT.

package document

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewRestColServiceCreateDocument2Params creates a new RestColServiceCreateDocument2Params object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRestColServiceCreateDocument2Params() *RestColServiceCreateDocument2Params {
	return &RestColServiceCreateDocument2Params{
		timeout: cr.DefaultTimeout,
	}
}

// NewRestColServiceCreateDocument2ParamsWithTimeout creates a new RestColServiceCreateDocument2Params object
// with the ability to set a timeout on a request.
func NewRestColServiceCreateDocument2ParamsWithTimeout(timeout time.Duration) *RestColServiceCreateDocument2Params {
	return &RestColServiceCreateDocument2Params{
		timeout: timeout,
	}
}

// NewRestColServiceCreateDocument2ParamsWithContext creates a new RestColServiceCreateDocument2Params object
// with the ability to set a context for a request.
func NewRestColServiceCreateDocument2ParamsWithContext(ctx context.Context) *RestColServiceCreateDocument2Params {
	return &RestColServiceCreateDocument2Params{
		Context: ctx,
	}
}

// NewRestColServiceCreateDocument2ParamsWithHTTPClient creates a new RestColServiceCreateDocument2Params object
// with the ability to set a custom HTTPClient for a request.
func NewRestColServiceCreateDocument2ParamsWithHTTPClient(client *http.Client) *RestColServiceCreateDocument2Params {
	return &RestColServiceCreateDocument2Params{
		HTTPClient: client,
	}
}

/*
RestColServiceCreateDocument2Params contains all the parameters to send to the API endpoint

	for the rest col service create document2 operation.

	Typically these are written to a http.Request.
*/
type RestColServiceCreateDocument2Params struct {

	// CollectionID.
	CollectionID string

	/* Data.

	   data represents rawdata for any kind of formating

	   Format: byte
	*/
	Data *strfmt.Base64

	// Dataformat.
	//
	// Default: "DATA_FORMAT_UNKNOWN"
	Dataformat *string

	// DocumentID.
	DocumentID *string

	// ProjectID.
	ProjectID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the rest col service create document2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceCreateDocument2Params) WithDefaults() *RestColServiceCreateDocument2Params {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the rest col service create document2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceCreateDocument2Params) SetDefaults() {
	var (
		dataformatDefault = string("DATA_FORMAT_UNKNOWN")
	)

	val := RestColServiceCreateDocument2Params{
		Dataformat: &dataformatDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithTimeout(timeout time.Duration) *RestColServiceCreateDocument2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithContext(ctx context.Context) *RestColServiceCreateDocument2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithHTTPClient(client *http.Client) *RestColServiceCreateDocument2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCollectionID adds the collectionID to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithCollectionID(collectionID string) *RestColServiceCreateDocument2Params {
	o.SetCollectionID(collectionID)
	return o
}

// SetCollectionID adds the collectionId to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetCollectionID(collectionID string) {
	o.CollectionID = collectionID
}

// WithData adds the data to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithData(data *strfmt.Base64) *RestColServiceCreateDocument2Params {
	o.SetData(data)
	return o
}

// SetData adds the data to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetData(data *strfmt.Base64) {
	o.Data = data
}

// WithDataformat adds the dataformat to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithDataformat(dataformat *string) *RestColServiceCreateDocument2Params {
	o.SetDataformat(dataformat)
	return o
}

// SetDataformat adds the dataformat to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetDataformat(dataformat *string) {
	o.Dataformat = dataformat
}

// WithDocumentID adds the documentID to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithDocumentID(documentID *string) *RestColServiceCreateDocument2Params {
	o.SetDocumentID(documentID)
	return o
}

// SetDocumentID adds the documentId to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetDocumentID(documentID *string) {
	o.DocumentID = documentID
}

// WithProjectID adds the projectID to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) WithProjectID(projectID *string) *RestColServiceCreateDocument2Params {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the rest col service create document2 params
func (o *RestColServiceCreateDocument2Params) SetProjectID(projectID *string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *RestColServiceCreateDocument2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param collectionId
	if err := r.SetPathParam("collectionId", o.CollectionID); err != nil {
		return err
	}

	if o.Data != nil {

		// query param data
		var qrData strfmt.Base64

		if o.Data != nil {
			qrData = *o.Data
		}
		qData := qrData.String()
		if qData != "" {

			if err := r.SetQueryParam("data", qData); err != nil {
				return err
			}
		}
	}

	if o.Dataformat != nil {

		// query param dataformat
		var qrDataformat string

		if o.Dataformat != nil {
			qrDataformat = *o.Dataformat
		}
		qDataformat := qrDataformat
		if qDataformat != "" {

			if err := r.SetQueryParam("dataformat", qDataformat); err != nil {
				return err
			}
		}
	}

	if o.DocumentID != nil {

		// query param documentId
		var qrDocumentID string

		if o.DocumentID != nil {
			qrDocumentID = *o.DocumentID
		}
		qDocumentID := qrDocumentID
		if qDocumentID != "" {

			if err := r.SetQueryParam("documentId", qDocumentID); err != nil {
				return err
			}
		}
	}

	if o.ProjectID != nil {

		// query param projectId
		var qrProjectID string

		if o.ProjectID != nil {
			qrProjectID = *o.ProjectID
		}
		qProjectID := qrProjectID
		if qProjectID != "" {

			if err := r.SetQueryParam("projectId", qProjectID); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
