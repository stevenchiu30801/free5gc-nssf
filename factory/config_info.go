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

type Info struct {

    Version string `yaml:"version,omitempty"`

    Description string `yaml:"description,omitempty"`

    Host string `yaml:"host"`

    ServingPlmnId *PlmnId `yaml:"servingPlmnId"`
}

func (i *Info) checkIntegrity() error {
    if i.Host == "" {
        return fmt.Errorf("`host` in configuration should not be empty")
    }

    if i.ServingPlmnId == nil {
        return fmt.Errorf("`servingPlmnId` in configuration should not be empty")
    } else {
        err := i.ServingPlmnId.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`servingPlmnId`:%s", err.Error())
        }
    }

    return nil
}
