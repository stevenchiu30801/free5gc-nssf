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

type MappingFromPlmnConfig struct {

    OperatorName string `yaml:"operatorName,omitempty"`

    HomePlmnId *PlmnId `yaml:"homePlmnId"`

    MappingOfSnssai []MappingOfSnssai `yaml:"mappingOfSnssai"`
}

func (m *MappingFromPlmnConfig) checkIntegrity() error {
    if m.HomePlmnId == nil {
        return fmt.Errorf("`homePlmnId` should not be empty")
    } else {
        err := m.HomePlmnId.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`homePlmnId`:%s", err.Error())
        }
    }

    if m.MappingOfSnssai == nil || len(m.MappingOfSnssai) == 0 {
        return fmt.Errorf("`mappingOfSnssai` should not empty")
    } else {
        for i, mappingOfSnssai := range m.MappingOfSnssai {
            err := mappingOfSnssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`mappingOfSnssai`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
