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

type TaConfig struct {
    Tac string `yaml:"tac"`
    SupportedNssai []Snssai `yaml:"supportedNssai"`
}

func (t *TaConfig) checkIntegrity() error {
    if t.Tac == "" {
        return errors.New("`tac` in configuration should not be empty")
    }

    if t.SupportedNssai == nil || len(t.SupportedNssai) == 0 {
        return errors.New("`supportedNssai` in configuration should not be empty")
    } else {
        for i, supportedSnssai := range t.SupportedNssai {
            err := supportedSnssai.CheckIntegrity()
            if err != nil {
                errMsg := "`supportedNssai`[" + strconv.Itoa(i) + "]:" + err.Error()
                return errors.New(errMsg)
            }
        }
    }

    return nil
}
