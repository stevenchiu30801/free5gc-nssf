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

type ConfiguredSnssai struct {

	ConfiguredSnssai *Snssai `json:"configuredSnssai"`

	MappedHomeSnssai *Snssai `json:"mappedHomeSnssai,omitempty"`
}

func (c *ConfiguredSnssai) CheckIntegrity() error {
    if c.ConfiguredSnssai == nil {
        return errors.New("`configuredSnssai` in query parameter should not be empty")
    } else {
        err := c.ConfiguredSnssai.CheckIntegrity()
        if err != nil {
            errMsg := "`configuredSnssai`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if c.MappedHomeSnssai != nil {
        err := c.MappedHomeSnssai.CheckIntegrity()
        if err != nil {
            errMsg := "`mappedHomeSnssai`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    return nil
}
