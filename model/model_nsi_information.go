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

type NsiInformation struct {

    NrfId string `json:"nrfId" yaml:"nrfId"`

    NsiId string `json:"nsiId,omitempty" yaml:"nsiId,omitempty"`
}

func (n *NsiInformation) CheckIntegrity() error {
    if n.NrfId == "" {
        return fmt.Errorf("`nrfId` should not be empty")
    }
    // TODO: Check whether `NrfId` is valid URI or not

    return nil
}
