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

    Service *Service `yaml:"service"`

    Version string `yaml:"version,omitempty"`

    Title string `yaml:"title,omitempty"`

    Description string `yaml:"description,omitempty"`

    Url string `yaml:"url"`

    ServingPlmnId *PlmnId `yaml:"servingPlmnId"`
}

func (i *Info) checkIntegrity() error {
    if i.Service == nil || *i.Service == Service("") {
        return fmt.Errorf("`service` in configuration should not be empty")
    } else {
        err := i.Service.checkIntegrity()
        if err != nil {
            return fmt.Errorf("`service`:%s", err.Error())
        }
    }

    if i.Url == "" {
        return fmt.Errorf("`url` in configuration should not be empty")
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
