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

type AmfConfig struct {

    NfId string `yaml:"nfId"`

    SupportedNssai []Snssai `yaml:"supportedNssai"`
}

func (a *AmfConfig) checkIntegrity() error {
    if a.NfId == "" {
        return errors.New("`nfId` in configuration should not be empty")
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
