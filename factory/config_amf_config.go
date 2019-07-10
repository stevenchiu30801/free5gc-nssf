/*
 * NSSF Configuration Factory
 */

package factory

import (
    "fmt"

    . "free5gc-nssf/model"
)

type AmfConfig struct {

    NfId string `yaml:"nfId"`

    SupportedNssaiAvailabilityData []SupportedNssaiAvailabilityData `yaml:"supportedNssaiAvailabilityData"`
}

func (a *AmfConfig) checkIntegrity() error {
    if a.NfId == "" {
        return fmt.Errorf("`nfId` should not be empty")
    }

    if a.SupportedNssaiAvailabilityData == nil || len(a.SupportedNssaiAvailabilityData) == 0 {
        return fmt.Errorf("`supportedNssaiAvailabilityData` should not be empty")
    } else {
        for i, supportedNssaiAvailabilityData := range a.SupportedNssaiAvailabilityData {
            err := supportedNssaiAvailabilityData.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedNssaiAvailabilityData`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
