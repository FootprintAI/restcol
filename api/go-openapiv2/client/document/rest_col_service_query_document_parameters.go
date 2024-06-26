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
	"github.com/go-openapi/swag"
)

// NewRestColServiceQueryDocumentParams creates a new RestColServiceQueryDocumentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewRestColServiceQueryDocumentParams() *RestColServiceQueryDocumentParams {
	return &RestColServiceQueryDocumentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewRestColServiceQueryDocumentParamsWithTimeout creates a new RestColServiceQueryDocumentParams object
// with the ability to set a timeout on a request.
func NewRestColServiceQueryDocumentParamsWithTimeout(timeout time.Duration) *RestColServiceQueryDocumentParams {
	return &RestColServiceQueryDocumentParams{
		timeout: timeout,
	}
}

// NewRestColServiceQueryDocumentParamsWithContext creates a new RestColServiceQueryDocumentParams object
// with the ability to set a context for a request.
func NewRestColServiceQueryDocumentParamsWithContext(ctx context.Context) *RestColServiceQueryDocumentParams {
	return &RestColServiceQueryDocumentParams{
		Context: ctx,
	}
}

// NewRestColServiceQueryDocumentParamsWithHTTPClient creates a new RestColServiceQueryDocumentParams object
// with the ability to set a custom HTTPClient for a request.
func NewRestColServiceQueryDocumentParamsWithHTTPClient(client *http.Client) *RestColServiceQueryDocumentParams {
	return &RestColServiceQueryDocumentParams{
		HTTPClient: client,
	}
}

/*
RestColServiceQueryDocumentParams contains all the parameters to send to the API endpoint

	for the rest col service query document operation.

	Typically these are written to a http.Request.
*/
type RestColServiceQueryDocumentParams struct {

	// CollectionID.
	CollectionID string

	/* EndedAt.

	   endedAt specifies when is the ended timeframe of the query

	   Format: date-time
	*/
	EndedAt *strfmt.DateTime

	/* FieldSelectors.

	   dot-concatenated fields, ex: fielda.fieldb.fieldc
	*/
	FieldSelectors []string

	// LimitCount.
	//
	// Format: int32
	LimitCount *int32

	// ProjectID.
	ProjectID string

	// SinceTs.
	//
	// Format: date-time
	SinceTs *strfmt.DateTime

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the rest col service query document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceQueryDocumentParams) WithDefaults() *RestColServiceQueryDocumentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the rest col service query document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *RestColServiceQueryDocumentParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithTimeout(timeout time.Duration) *RestColServiceQueryDocumentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithContext(ctx context.Context) *RestColServiceQueryDocumentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithHTTPClient(client *http.Client) *RestColServiceQueryDocumentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCollectionID adds the collectionID to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithCollectionID(collectionID string) *RestColServiceQueryDocumentParams {
	o.SetCollectionID(collectionID)
	return o
}

// SetCollectionID adds the collectionId to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetCollectionID(collectionID string) {
	o.CollectionID = collectionID
}

// WithEndedAt adds the endedAt to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithEndedAt(endedAt *strfmt.DateTime) *RestColServiceQueryDocumentParams {
	o.SetEndedAt(endedAt)
	return o
}

// SetEndedAt adds the endedAt to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetEndedAt(endedAt *strfmt.DateTime) {
	o.EndedAt = endedAt
}

// WithFieldSelectors adds the fieldSelectors to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithFieldSelectors(fieldSelectors []string) *RestColServiceQueryDocumentParams {
	o.SetFieldSelectors(fieldSelectors)
	return o
}

// SetFieldSelectors adds the fieldSelectors to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetFieldSelectors(fieldSelectors []string) {
	o.FieldSelectors = fieldSelectors
}

// WithLimitCount adds the limitCount to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithLimitCount(limitCount *int32) *RestColServiceQueryDocumentParams {
	o.SetLimitCount(limitCount)
	return o
}

// SetLimitCount adds the limitCount to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetLimitCount(limitCount *int32) {
	o.LimitCount = limitCount
}

// WithProjectID adds the projectID to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithProjectID(projectID string) *RestColServiceQueryDocumentParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WithSinceTs adds the sinceTs to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) WithSinceTs(sinceTs *strfmt.DateTime) *RestColServiceQueryDocumentParams {
	o.SetSinceTs(sinceTs)
	return o
}

// SetSinceTs adds the sinceTs to the rest col service query document params
func (o *RestColServiceQueryDocumentParams) SetSinceTs(sinceTs *strfmt.DateTime) {
	o.SinceTs = sinceTs
}

// WriteToRequest writes these params to a swagger request
func (o *RestColServiceQueryDocumentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param collectionId
	if err := r.SetPathParam("collectionId", o.CollectionID); err != nil {
		return err
	}

	if o.EndedAt != nil {

		// query param endedAt
		var qrEndedAt strfmt.DateTime

		if o.EndedAt != nil {
			qrEndedAt = *o.EndedAt
		}
		qEndedAt := qrEndedAt.String()
		if qEndedAt != "" {

			if err := r.SetQueryParam("endedAt", qEndedAt); err != nil {
				return err
			}
		}
	}

	if o.FieldSelectors != nil {

		// binding items for fieldSelectors
		joinedFieldSelectors := o.bindParamFieldSelectors(reg)

		// query array param fieldSelectors
		if err := r.SetQueryParam("fieldSelectors", joinedFieldSelectors...); err != nil {
			return err
		}
	}

	if o.LimitCount != nil {

		// query param limitCount
		var qrLimitCount int32

		if o.LimitCount != nil {
			qrLimitCount = *o.LimitCount
		}
		qLimitCount := swag.FormatInt32(qrLimitCount)
		if qLimitCount != "" {

			if err := r.SetQueryParam("limitCount", qLimitCount); err != nil {
				return err
			}
		}
	}

	// path param projectId
	if err := r.SetPathParam("projectId", o.ProjectID); err != nil {
		return err
	}

	if o.SinceTs != nil {

		// query param sinceTs
		var qrSinceTs strfmt.DateTime

		if o.SinceTs != nil {
			qrSinceTs = *o.SinceTs
		}
		qSinceTs := qrSinceTs.String()
		if qSinceTs != "" {

			if err := r.SetQueryParam("sinceTs", qSinceTs); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindParamRestColServiceQueryDocument binds the parameter fieldSelectors
func (o *RestColServiceQueryDocumentParams) bindParamFieldSelectors(formats strfmt.Registry) []string {
	fieldSelectorsIR := o.FieldSelectors

	var fieldSelectorsIC []string
	for _, fieldSelectorsIIR := range fieldSelectorsIR { // explode []string

		fieldSelectorsIIV := fieldSelectorsIIR // string as string
		fieldSelectorsIC = append(fieldSelectorsIC, fieldSelectorsIIV)
	}

	// items.CollectionFormat: "multi"
	fieldSelectorsIS := swag.JoinByFormat(fieldSelectorsIC, "multi")

	return fieldSelectorsIS
}
