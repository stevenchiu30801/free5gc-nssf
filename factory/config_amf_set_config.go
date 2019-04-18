/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"

    . "../model"
)

type AmfSetConfig struct {

    AmfSetId string `yaml:"amfSetId"`

    AmfList []string `yaml:"amfList,omitempty"`

    SupportedNssai []Snssai `yaml:"supportedNssai"`
}

func (a *AmfSetConfig) checkIntegrity() error {
    if a.AmfSetId == "" {
        return fmt.Errorf("`amfSetId` in configuration should not be empty")
    }

    if a.SupportedNssai == nil || len(a.SupportedNssai) == 0 {
        return fmt.Errorf("`supportedNssai` in configuration should not be empty")
    } else {
        for i, supportedSnssai := range a.SupportedNssai {
            err := supportedSnssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedNssai`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
