/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"
)

type Service string

// List of NSSF service type
const (
    NSSF_NSSELECTION Service = "Nnssf-NSSelection"
    NSSF_NSSAIAVAILABILITY Service = "Nnssf-NSSAIAvailability"
)

func (s *Service) checkIntegrity() error {
    if *s != NSSF_NSSELECTION && *s != NSSF_NSSAIAVAILABILITY {
        return fmt.Errorf("'%s' is unrecognized", string(*s))
    }

    return nil
}
