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

type TaConfig struct {

    Tac string `yaml:"tac"`

    AccessType *AccessType `yaml:"accessType"`

    SupportedNssai []Snssai `yaml:"supportedNssai"`
}

func (t *TaConfig) checkIntegrity() error {
    if t.Tac == "" {
        return fmt.Errorf("`tac` in configuration should not be empty")
    }

    if t.AccessType == nil || *t.AccessType == AccessType("") {
        return fmt.Errorf("`accessType` in configuration should not be empty")
    } else {
        err := t.AccessType.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`accessType`:%s", err.Error())
        }
    }

    if t.SupportedNssai == nil || len(t.SupportedNssai) == 0 {
        return fmt.Errorf("`supportedNssai` in configuration should not be empty")
    } else {
        for i, supportedSnssai := range t.SupportedNssai {
            err := supportedSnssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedNssai`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
