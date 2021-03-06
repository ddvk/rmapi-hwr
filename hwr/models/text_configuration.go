// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TextConfiguration TextConfiguration
//
// swagger:model TextConfiguration
type TextConfiguration struct {

	// configuration
	Configuration *TextConfConfiguration `json:"configuration,omitempty"`

	// guides
	Guides *GuidesConfiguration `json:"guides,omitempty"`

	// margin
	Margin *MarginConfiguration `json:"margin,omitempty"`
}

// Validate validates this text configuration
func (m *TextConfiguration) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfiguration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGuides(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMargin(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TextConfiguration) validateConfiguration(formats strfmt.Registry) error {

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

func (m *TextConfiguration) validateGuides(formats strfmt.Registry) error {

	if swag.IsZero(m.Guides) { // not required
		return nil
	}

	if m.Guides != nil {
		if err := m.Guides.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("guides")
			}
			return err
		}
	}

	return nil
}

func (m *TextConfiguration) validateMargin(formats strfmt.Registry) error {

	if swag.IsZero(m.Margin) { // not required
		return nil
	}

	if m.Margin != nil {
		if err := m.Margin.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("margin")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TextConfiguration) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TextConfiguration) UnmarshalBinary(b []byte) error {
	var res TextConfiguration
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
