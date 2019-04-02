/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

import (
    "errors"
)

type SliceInfoForPduSession struct {

	SNssai *Snssai `json:"sNssai"`

	RoamingIndication *RoamingIndication `json:"roamingIndication"`

	HomeSnssai *Snssai `json:"homeSnssai,omitempty"`
}

func (s *SliceInfoForPduSession) CheckIntegrity() error {
    if s.SNssai == nil {
        return errors.New("`sNssai` in query parameter should not be empty")
    }

    if s.RoamingIndication == nil {
        return errors.New("`roamingIndication` in query parameter should not be empty")
    }

    return nil
}
