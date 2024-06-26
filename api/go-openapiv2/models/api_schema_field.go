// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// APISchemaField api schema field
//
// swagger:model apiSchemaField
type APISchemaField struct {

	// datatype
	Datatype *APISchemaFieldDataType `json:"datatype,omitempty"`

	// example
	Example interface{} `json:"example,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this api schema field
func (m *APISchemaField) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDatatype(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APISchemaField) validateDatatype(formats strfmt.Registry) error {
	if swag.IsZero(m.Datatype) { // not required
		return nil
	}

	if m.Datatype != nil {
		if err := m.Datatype.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datatype")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datatype")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this api schema field based on the context it is used
func (m *APISchemaField) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDatatype(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APISchemaField) contextValidateDatatype(ctx context.Context, formats strfmt.Registry) error {

	if m.Datatype != nil {

		if swag.IsZero(m.Datatype) { // not required
			return nil
		}

		if err := m.Datatype.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("datatype")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("datatype")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *APISchemaField) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APISchemaField) UnmarshalBinary(b []byte) error {
	var res APISchemaField
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
