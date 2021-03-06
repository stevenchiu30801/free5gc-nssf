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
    "fmt"
)

type PlmnId struct {

    Mcc string `json:"mcc" yaml:"mcc"`

    Mnc string `json:"mnc" yaml:"mnc"`
}

func (p *PlmnId) CheckIntegrity() error {
    if p.Mcc == "" {
        return fmt.Errorf("`mcc` should not be empty")
    }

    if p.Mnc == "" {
        return fmt.Errorf("`mnc` should not be empty")
    }

    return nil
}
