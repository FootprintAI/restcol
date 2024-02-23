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

// APICollectionType api collection type
//
// swagger:model apiCollectionType
type APICollectionType string

func NewAPICollectionType(value APICollectionType) *APICollectionType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated APICollectionType.
func (m APICollectionType) Pointer() *APICollectionType {
	return &m
}

const (

	// APICollectionTypeCOLLECTIONTYPENONE captures enum value "COLLECTION_TYPE_NONE"
	APICollectionTypeCOLLECTIONTYPENONE APICollectionType = "COLLECTION_TYPE_NONE"

	// APICollectionTypeCOLLECTIONTYPEREGULARFILES captures enum value "COLLECTION_TYPE_REGULAR_FILES"
	APICollectionTypeCOLLECTIONTYPEREGULARFILES APICollectionType = "COLLECTION_TYPE_REGULAR_FILES"

	// APICollectionTypeCOLLECTIONTYPETIMESERIES captures enum value "COLLECTION_TYPE_TIMESERIES"
	APICollectionTypeCOLLECTIONTYPETIMESERIES APICollectionType = "COLLECTION_TYPE_TIMESERIES"

	// APICollectionTypeCOLLECTIONTYPETRANSACTION captures enum value "COLLECTION_TYPE_TRANSACTION"
	APICollectionTypeCOLLECTIONTYPETRANSACTION APICollectionType = "COLLECTION_TYPE_TRANSACTION"

	// APICollectionTypeCOLLECTIONTYPEVECTOR captures enum value "COLLECTION_TYPE_VECTOR"
	APICollectionTypeCOLLECTIONTYPEVECTOR APICollectionType = "COLLECTION_TYPE_VECTOR"

	// APICollectionTypeCOLLECTIONTYPEPROXY captures enum value "COLLECTION_TYPE_PROXY"
	APICollectionTypeCOLLECTIONTYPEPROXY APICollectionType = "COLLECTION_TYPE_PROXY"
)

// for schema
var apiCollectionTypeEnum []interface{}

func init() {
	var res []APICollectionType
	if err := json.Unmarshal([]byte(`["COLLECTION_TYPE_NONE","COLLECTION_TYPE_REGULAR_FILES","COLLECTION_TYPE_TIMESERIES","COLLECTION_TYPE_TRANSACTION","COLLECTION_TYPE_VECTOR","COLLECTION_TYPE_PROXY"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		apiCollectionTypeEnum = append(apiCollectionTypeEnum, v)
	}
}

func (m APICollectionType) validateAPICollectionTypeEnum(path, location string, value APICollectionType) error {
	if err := validate.EnumCase(path, location, value, apiCollectionTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this api collection type
func (m APICollectionType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateAPICollectionTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this api collection type based on context it is used
func (m APICollectionType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
