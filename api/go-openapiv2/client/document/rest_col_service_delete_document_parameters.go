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

// NewRestColServiceDeleteDocumentParams creates a new RestColServiceDeleteDocumentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRestColServiceDeleteDocumentParams() *RestColServiceDeleteDocumentParams {
	return &RestColServiceDeleteDocumentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewRestColServiceDeleteDocumentParamsWithTimeout creates a new RestColServiceDeleteDocumentParams object
// with the ability to set a timeout on a request.
func NewRestColServiceDeleteDocumentParamsWithTimeout(timeout time.Duration) *RestColServiceDeleteDocumentParams {
	return &RestColServiceDeleteDocumentParams{
		timeout: timeout,
	}
}

// NewRestColServiceDeleteDocumentParamsWithContext creates a new RestColServiceDeleteDocumentParams object
// with the ability to set a context for a request.
func NewRestColServiceDeleteDocumentParamsWithContext(ctx context.Context) *RestColServiceDeleteDocumentParams {
	return &RestColServiceDeleteDocumentParams{
		Context: ctx,
	}
}

// NewRestColServiceDeleteDocumentParamsWithHTTPClient creates a new RestColServiceDeleteDocumentParams object
// with the ability to set a custom HTTPClient for a request.
func NewRestColServiceDeleteDocumentParamsWithHTTPClient(client *http.Client) *RestColServiceDeleteDocumentParams {
	return &RestColServiceDeleteDocumentParams{
		HTTPClient: client,
	}
}

/*
RestColServiceDeleteDocumentParams contains all the parameters to send to the API endpoint

	for the rest col service delete document operation.

	Typically these are written to a http.Request.
*/
type RestColServiceDeleteDocumentParams struct {

	// Cid.
	Cid string

	// Did.
	Did string

	// Pid.
	Pid *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the rest col service delete document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceDeleteDocumentParams) WithDefaults() *RestColServiceDeleteDocumentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the rest col service delete document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceDeleteDocumentParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithTimeout(timeout time.Duration) *RestColServiceDeleteDocumentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithContext(ctx context.Context) *RestColServiceDeleteDocumentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithHTTPClient(client *http.Client) *RestColServiceDeleteDocumentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCid adds the cid to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithCid(cid string) *RestColServiceDeleteDocumentParams {
	o.SetCid(cid)
	return o
}

// SetCid adds the cid to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetCid(cid string) {
	o.Cid = cid
}

// WithDid adds the did to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithDid(did string) *RestColServiceDeleteDocumentParams {
	o.SetDid(did)
	return o
}

// SetDid adds the did to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetDid(did string) {
	o.Did = did
}

// WithPid adds the pid to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) WithPid(pid *string) *RestColServiceDeleteDocumentParams {
	o.SetPid(pid)
	return o
}

// SetPid adds the pid to the rest col service delete document params
func (o *RestColServiceDeleteDocumentParams) SetPid(pid *string) {
	o.Pid = pid
}

// WriteToRequest writes these params to a swagger request
func (o *RestColServiceDeleteDocumentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cid
	if err := r.SetPathParam("cid", o.Cid); err != nil {
		return err
	}

	// path param did
	if err := r.SetPathParam("did", o.Did); err != nil {
		return err
	}

	if o.Pid != nil {

		// query param pid
		var qrPid string

		if o.Pid != nil {
			qrPid = *o.Pid
		}
		qPid := qrPid
		if qPid != "" {

			if err := r.SetQueryParam("pid", qPid); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
