// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BatchInput BatchInput
//
// swagger:model BatchInput
type BatchInput struct {

	// The configuration for the recognition.
	Configuration *Configuration `json:"configuration,omitempty"`

	// recognition type
	// Required: true
	// Enum: [Text Math Diagram Raw Content Text Document]
	ContentType *string `json:"contentType"`

	// target of conversion, no conversion will be made if that parameter is not provided
	// Enum: [DIGITAL_PUBLISH DIGITAL_EDIT]
	ConversionState string `json:"conversionState,omitempty"`

	// height of the writing area
	Height int32 `json:"height,omitempty"`

	// The write entries that corresponds to the input iink
	// Required: true
	StrokeGroups []*StrokeGroup `json:"strokeGroups"`

	// A global CSS styling for your content. See https://developer.myscript.com/docs/interactive-ink/latestweb/myscriptjs/styling/
	Theme string `json:"theme,omitempty"`

	// width of the writing area
	Width int32 `json:"width,omitempty"`

	// x resolution of the writing area in dpi
	XDPI float32 `json:"xDPI,omitempty"`

	// y resolution of the writing area in dpi
	YDPI float32 `json:"yDPI,omitempty"`
}

// Validate validates this batch input
func (m *BatchInput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfiguration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateContentType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConversionState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStrokeGroups(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BatchInput) validateConfiguration(formats strfmt.Registry) error {

	if swag.IsZero(m.Configuration) { // not required
		return nil
	}

	if m.Configuration != nil {
		if err := m.Configuration.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("configuration")
			}
			return err
		}
	}

	return nil
}

var batchInputTypeContentTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Text","Math","Diagram","Raw Content","Text Document"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		batchInputTypeContentTypePropEnum = append(batchInputTypeContentTypePropEnum, v)
	}
}

const (

	// BatchInputContentTypeText captures enum value "Text"
	BatchInputContentTypeText string = "Text"

	// BatchInputContentTypeMath captures enum value "Math"
	BatchInputContentTypeMath string = "Math"

	// BatchInputContentTypeDiagram captures enum value "Diagram"
	BatchInputContentTypeDiagram string = "Diagram"

	// BatchInputContentTypeRawContent captures enum value "Raw Content"
	BatchInputContentTypeRawContent string = "Raw Content"

	// BatchInputContentTypeTextDocument captures enum value "Text Document"
	BatchInputContentTypeTextDocument string = "Text Document"
)

// prop value enum
func (m *BatchInput) validateContentTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, batchInputTypeContentTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *BatchInput) validateContentType(formats strfmt.Registry) error {

	if err := validate.Required("contentType", "body", m.ContentType); err != nil {
		return err
	}

	// value enum
	if err := m.validateContentTypeEnum("contentType", "body", *m.ContentType); err != nil {
		return err
	}

	return nil
}

var batchInputTypeConversionStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["DIGITAL_PUBLISH","DIGITAL_EDIT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		batchInputTypeConversionStatePropEnum = append(batchInputTypeConversionStatePropEnum, v)
	}
}

const (

	// BatchInputConversionStateDIGITALPUBLISH captures enum value "DIGITAL_PUBLISH"
	BatchInputConversionStateDIGITALPUBLISH string = "DIGITAL_PUBLISH"

	// BatchInputConversionStateDIGITALEDIT captures enum value "DIGITAL_EDIT"
	BatchInputConversionStateDIGITALEDIT string = "DIGITAL_EDIT"
)

// prop value enum
func (m *BatchInput) validateConversionStateEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, batchInputTypeConversionStatePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *BatchInput) validateConversionState(formats strfmt.Registry) error {

	if swag.IsZero(m.ConversionState) { // not required
		return nil
	}

	// value enum
	if err := m.validateConversionStateEnum("conversionState", "body", m.ConversionState); err != nil {
		return err
	}

	return nil
}

func (m *BatchInput) validateStrokeGroups(formats strfmt.Registry) error {

	if err := validate.Required("strokeGroups", "body", m.StrokeGroups); err != nil {
		return err
	}

	for i := 0; i < len(m.StrokeGroups); i++ {
		if swag.IsZero(m.StrokeGroups[i]) { // not required
			continue
		}

		if m.StrokeGroups[i] != nil {
			if err := m.StrokeGroups[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("strokeGroups" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BatchInput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BatchInput) UnmarshalBinary(b []byte) error {
	var res BatchInput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
