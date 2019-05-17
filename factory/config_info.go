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

    ServingPlmnIdList []PlmnId `yaml:"servingPlmnIdList"`
}

func (i *Info) checkIntegrity() error {
    if i.Host == "" {
        return fmt.Errorf("`host` should not be empty")
    }

    if i.ServingPlmnIdList == nil || len(i.ServingPlmnIdList) == 0 {
        return fmt.Errorf("`servingPlmnIdList` should not be empty")
    } else {
        for i, plmnId := range i.ServingPlmnIdList {
            err := plmnId.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`servingPlmnIdList`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
