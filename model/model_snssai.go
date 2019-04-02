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

type Snssai struct {

    Sst int32 `json:"sst" yaml:"sst"`

    Sd string `json:"sd,omitempty" yaml:"sd,omitempty"`
}

func (s *Snssai) CheckIntegrity() error {
    if s.Sst == 0 {
        return errors.New("`sst` in query parameter should not be empty")
    }

    return nil
}
