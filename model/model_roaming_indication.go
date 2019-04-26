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

type RoamingIndication string

// List of RoamingIndication
const (
	NON_ROAMING RoamingIndication = "NON_ROAMING"
	LOCAL_BREAKOUT RoamingIndication = "LOCAL_BREAKOUT"
	HOME_ROUTED_ROAMING RoamingIndication = "HOME_ROUTED_ROAMING"
)

func (r *RoamingIndication) CheckIntegrity() error {
    if *r != NON_ROAMING && *r != LOCAL_BREAKOUT && *r != HOME_ROUTED_ROAMING {
        return fmt.Errorf("'%s' is unrecognized", string(*r))
    }

    return nil
}
