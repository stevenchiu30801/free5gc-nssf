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

    Tai *Tai `yaml:"tai"`

    AccessType *AccessType `yaml:"accessType"`

    SupportedSnssaiList []Snssai `yaml:"supportedSnssaiList,omitempty"`
}

func (t *TaConfig) checkIntegrity() error {
    if t.Tai == nil {
        return fmt.Errorf("`tac` in configuration should not be empty")
    } else {
        err := t.Tai.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`tai`:%s", err.Error())
        }
    }

    if t.AccessType == nil || *t.AccessType == AccessType("") {
        return fmt.Errorf("`accessType` in configuration should not be empty")
    } else {
        err := t.AccessType.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`accessType`:%s", err.Error())
        }
    }

    if t.SupportedSnssaiList != nil && len(t.SupportedSnssaiList) != 0 {
        for i, snssai := range t.SupportedSnssaiList {
            err := snssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedSnssaiList`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
