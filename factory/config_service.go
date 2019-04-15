/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "errors"
)

type Service string

// List of NSSF service type
const (
    NSSF_NSSELECTION Service = "Nnssf_NSSelection"
    NSSF_NSSAIAVAILABILITY Service = "Nnssf_NSSAIAvailability"
)

func (s *Service) checkIntegrity() error {
    if *s != NSSF_NSSELECTION && *s != NSSF_NSSAIAVAILABILITY {
        errMsg := "'" + string(*s) + "' is unrecognized"
        return errors.New(errMsg)
    }

    return nil
}
