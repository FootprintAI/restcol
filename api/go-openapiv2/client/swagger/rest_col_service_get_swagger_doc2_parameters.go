// Code generated by go-swagger; DO NOT EDIT.

package swagger

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

// NewRestColServiceGetSwaggerDoc2Params creates a new RestColServiceGetSwaggerDoc2Params object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRestColServiceGetSwaggerDoc2Params() *RestColServiceGetSwaggerDoc2Params {
	return &RestColServiceGetSwaggerDoc2Params{
		timeout: cr.DefaultTimeout,
	}
}

// NewRestColServiceGetSwaggerDoc2ParamsWithTimeout creates a new RestColServiceGetSwaggerDoc2Params object
// with the ability to set a timeout on a request.
func NewRestColServiceGetSwaggerDoc2ParamsWithTimeout(timeout time.Duration) *RestColServiceGetSwaggerDoc2Params {
	return &RestColServiceGetSwaggerDoc2Params{
		timeout: timeout,
	}
}

// NewRestColServiceGetSwaggerDoc2ParamsWithContext creates a new RestColServiceGetSwaggerDoc2Params object
// with the ability to set a context for a request.
func NewRestColServiceGetSwaggerDoc2ParamsWithContext(ctx context.Context) *RestColServiceGetSwaggerDoc2Params {
	return &RestColServiceGetSwaggerDoc2Params{
		Context: ctx,
	}
}

// NewRestColServiceGetSwaggerDoc2ParamsWithHTTPClient creates a new RestColServiceGetSwaggerDoc2Params object
// with the ability to set a custom HTTPClient for a request.
func NewRestColServiceGetSwaggerDoc2ParamsWithHTTPClient(client *http.Client) *RestColServiceGetSwaggerDoc2Params {
	return &RestColServiceGetSwaggerDoc2Params{
		HTTPClient: client,
	}
}

/*
RestColServiceGetSwaggerDoc2Params contains all the parameters to send to the API endpoint

	for the rest col service get swagger doc2 operation.

	Typically these are written to a http.Request.
*/
type RestColServiceGetSwaggerDoc2Params struct {

	// Pid.
	Pid string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the rest col service get swagger doc2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceGetSwaggerDoc2Params) WithDefaults() *RestColServiceGetSwaggerDoc2Params {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the rest col service get swagger doc2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceGetSwaggerDoc2Params) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) WithTimeout(timeout time.Duration) *RestColServiceGetSwaggerDoc2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) WithContext(ctx context.Context) *RestColServiceGetSwaggerDoc2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) WithHTTPClient(client *http.Client) *RestColServiceGetSwaggerDoc2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPid adds the pid to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) WithPid(pid string) *RestColServiceGetSwaggerDoc2Params {
	o.SetPid(pid)
	return o
}

// SetPid adds the pid to the rest col service get swagger doc2 params
func (o *RestColServiceGetSwaggerDoc2Params) SetPid(pid string) {
	o.Pid = pid
}

// WriteToRequest writes these params to a swagger request
func (o *RestColServiceGetSwaggerDoc2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param pid
	if err := r.SetPathParam("pid", o.Pid); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
