/*
 * NSSF Configuration Factory
 */

package factory

import (
    "fmt"

    . "free5gc-nssf/model"
)

type AmfSetConfig struct {

    AmfSetId string `yaml:"amfSetId"`

    AmfList []string `yaml:"amfList,omitempty"`

    NrfAmfSet string `yaml:"nrfAmfSet,omitempty"`

    SupportedNssaiAvailabilityData []SupportedNssaiAvailabilityData `yaml:"supportedNssaiAvailabilityData"`
}

func (a *AmfSetConfig) checkIntegrity() error {
    if a.AmfSetId == "" {
        return fmt.Errorf("`amfSetId` should not be empty")
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
