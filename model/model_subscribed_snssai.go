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

type SubscribedSnssai struct {

	SubscribedSnssai *Snssai `json:"subscribedSnssai"`

	DefaultIndication bool `json:"defaultIndication,omitempty"`
}

func (s *SubscribedSnssai) CheckIntegrity() error {
    if s.SubscribedSnssai == nil {
        return fmt.Errorf("`subscribedSnssai` should not be empty")
    } else {
        err := s.SubscribedSnssai.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`subscribedSnssai`:%s", err.Error())
        }
    }

    return nil
}
