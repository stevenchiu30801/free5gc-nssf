/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "errors"
    "strconv"

    . "../model"
)

type AmfSetConfig struct {

    AmfSetId string `yaml:"amfSetId"`

    AmfList []string `yaml:"amfList,omitempty"`

    SupportedNssai []Snssai `yaml:"supportedNssai"`
}

func (a *AmfSetConfig) checkIntegrity() error {
    if a.AmfSetId == "" {
        return errors.New("`amfSetId` in configuration should not be empty")
    }

    if a.SupportedNssai == nil || len(a.SupportedNssai) == 0 {
        return errors.New("`supportedNssai` in configuration should not be empty")
    } else {
        for i, supportedSnssai := range a.SupportedNssai {
            err := supportedSnssai.CheckIntegrity()
            if err != nil {
                errMsg := "`supportedNssai`[" + strconv.Itoa(i) + "]:" + err.Error()
                return errors.New(errMsg)
            }
        }
    }

    return nil
}
