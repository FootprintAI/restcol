// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// APIDataFormat api data format
//
// swagger:model apiDataFormat
type APIDataFormat string

func NewAPIDataFormat(value APIDataFormat) *APIDataFormat {
	return &value
}

// Pointer returns a pointer to a freshly-allocated APIDataFormat.
func (m APIDataFormat) Pointer() *APIDataFormat {
	return &m
}

const (

	// APIDataFormatDATAFORMATAUTO captures enum value "DATA_FORMAT_AUTO"
	APIDataFormatDATAFORMATAUTO APIDataFormat = "DATA_FORMAT_AUTO"

	// APIDataFormatDATAFORMATJSON captures enum value "DATA_FORMAT_JSON"
	APIDataFormatDATAFORMATJSON APIDataFormat = "DATA_FORMAT_JSON"

	// APIDataFormatDATAFORMATCSV captures enum value "DATA_FORMAT_CSV"
	APIDataFormatDATAFORMATCSV APIDataFormat = "DATA_FORMAT_CSV"

	// APIDataFormatDATAFORMATXML captures enum value "DATA_FORMAT_XML"
	APIDataFormatDATAFORMATXML APIDataFormat = "DATA_FORMAT_XML"

	// APIDataFormatDATAFORMATURL captures enum value "DATA_FORMAT_URL"
	APIDataFormatDATAFORMATURL APIDataFormat = "DATA_FORMAT_URL"

	// APIDataFormatDATAFORMATMEDIA captures enum value "DATA_FORMAT_MEDIA"
	APIDataFormatDATAFORMATMEDIA APIDataFormat = "DATA_FORMAT_MEDIA"
)

// for schema
var apiDataFormatEnum []interface{}

func init() {
	var res []APIDataFormat
	if err := json.Unmarshal([]byte(`["DATA_FORMAT_AUTO","DATA_FORMAT_JSON","DATA_FORMAT_CSV","DATA_FORMAT_XML","DATA_FORMAT_URL","DATA_FORMAT_MEDIA"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		apiDataFormatEnum = append(apiDataFormatEnum, v)
	}
}

func (m APIDataFormat) validateAPIDataFormatEnum(path, location string, value APIDataFormat) error {
	if err := validate.EnumCase(path, location, value, apiDataFormatEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this api data format
func (m APIDataFormat) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateAPIDataFormatEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this api data format based on context it is used
func (m APIDataFormat) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
