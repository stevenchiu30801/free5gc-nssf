/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "errors"

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
        return errors.New("`service` in configuration should not be empty")
    } else {
        err := i.Service.checkIntegrity()
        if err != nil {
            errMsg := "`service`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if i.Url == "" {
        return errors.New("`url` in configuration should not be empty")
    }

    if i.ServingPlmnId == nil {
        return errors.New("`servingPlmnId` in configuration should not be empty")
    } else {
        err := i.ServingPlmnId.CheckIntegrity()
        if err != nil {
            errMsg := "`servingPlmnId`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    return nil
}
