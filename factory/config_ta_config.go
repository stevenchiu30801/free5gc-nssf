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

    SupportedSnssaiList []Snssai `yaml:"supportedSnssaiList"`

    RestrictedSnssaiList []RestrictedSnssai `yaml:"restrictedSnssaiList,omitempty"`
}

func (t *TaConfig) checkIntegrity() error {
    if t.Tai == nil {
        return fmt.Errorf("`tac` should not be empty")
    } else {
        err := t.Tai.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`tai`:%s", err.Error())
        }
    }

    if t.AccessType == nil || *t.AccessType == AccessType("") {
        return fmt.Errorf("`accessType` should not be empty")
    } else {
        err := t.AccessType.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`accessType`:%s", err.Error())
        }
    }

    if t.SupportedSnssaiList == nil || len(t.SupportedSnssaiList) == 0 {
        return fmt.Errorf("`supportedSnssaiList` should not be empty")
    } else {
        for i, snssai := range t.SupportedSnssaiList {
            err := snssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedSnssaiList`[%d]:%s", i, err.Error())
            }
        }
    }

    if t.RestrictedSnssaiList != nil && len(t.RestrictedSnssaiList) != 0 {
        for i, restrictedSnssai := range t.RestrictedSnssaiList {
            err := restrictedSnssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`restrictedSnssaiList`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
