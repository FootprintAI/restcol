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

// NewRestColServiceGetDocument2Params creates a new RestColServiceGetDocument2Params object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRestColServiceGetDocument2Params() *RestColServiceGetDocument2Params {
	return &RestColServiceGetDocument2Params{
		timeout: cr.DefaultTimeout,
	}
}

// NewRestColServiceGetDocument2ParamsWithTimeout creates a new RestColServiceGetDocument2Params object
// with the ability to set a timeout on a request.
func NewRestColServiceGetDocument2ParamsWithTimeout(timeout time.Duration) *RestColServiceGetDocument2Params {
	return &RestColServiceGetDocument2Params{
		timeout: timeout,
	}
}

// NewRestColServiceGetDocument2ParamsWithContext creates a new RestColServiceGetDocument2Params object
// with the ability to set a context for a request.
func NewRestColServiceGetDocument2ParamsWithContext(ctx context.Context) *RestColServiceGetDocument2Params {
	return &RestColServiceGetDocument2Params{
		Context: ctx,
	}
}

// NewRestColServiceGetDocument2ParamsWithHTTPClient creates a new RestColServiceGetDocument2Params object
// with the ability to set a custom HTTPClient for a request.
func NewRestColServiceGetDocument2ParamsWithHTTPClient(client *http.Client) *RestColServiceGetDocument2Params {
	return &RestColServiceGetDocument2Params{
		HTTPClient: client,
	}
}

/*
RestColServiceGetDocument2Params contains all the parameters to send to the API endpoint

	for the rest col service get document2 operation.

	Typically these are written to a http.Request.
*/
type RestColServiceGetDocument2Params struct {

	// Cid.
	Cid string

	// Did.
	Did string

	// Pid.
	Pid string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the rest col service get document2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceGetDocument2Params) WithDefaults() *RestColServiceGetDocument2Params {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the rest col service get document2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceGetDocument2Params) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithTimeout(timeout time.Duration) *RestColServiceGetDocument2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithContext(ctx context.Context) *RestColServiceGetDocument2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithHTTPClient(client *http.Client) *RestColServiceGetDocument2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCid adds the cid to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithCid(cid string) *RestColServiceGetDocument2Params {
	o.SetCid(cid)
	return o
}

// SetCid adds the cid to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetCid(cid string) {
	o.Cid = cid
}

// WithDid adds the did to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithDid(did string) *RestColServiceGetDocument2Params {
	o.SetDid(did)
	return o
}

// SetDid adds the did to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetDid(did string) {
	o.Did = did
}

// WithPid adds the pid to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) WithPid(pid string) *RestColServiceGetDocument2Params {
	o.SetPid(pid)
	return o
}

// SetPid adds the pid to the rest col service get document2 params
func (o *RestColServiceGetDocument2Params) SetPid(pid string) {
	o.Pid = pid
}

// WriteToRequest writes these params to a swagger request
func (o *RestColServiceGetDocument2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param pid
	if err := r.SetPathParam("pid", o.Pid); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
