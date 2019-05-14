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

    return nil
}
