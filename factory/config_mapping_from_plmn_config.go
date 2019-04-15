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

type MappingFromPlmnConfig struct {
    OperatorName string `yaml:"operatorName,omitempty"`
    HomePlmnId *PlmnId `yaml:"homePlmnId"`
    MappingOfSnssai []MappingOfSnssai `yaml:"mappingOfSnssai"`
}

func (m *MappingFromPlmnConfig) checkIntegrity() error {
    if m.HomePlmnId == nil {
        return errors.New("`homePlmnId` in configuration should not be empty")
    } else {
        err := m.HomePlmnId.CheckIntegrity()
        if err != nil {
            errMsg := "`homePlmnId`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if m.MappingOfSnssai == nil || len(m.MappingOfSnssai) == 0 {
        return errors.New("`mappingOfSnssai` in configuration should not empty")
    } else {
        for i, mappingOfSnssai := range m.MappingOfSnssai {
            err := mappingOfSnssai.CheckIntegrity()
            if err != nil {
                errMsg := "`mappingOfSnssai`[" + strconv.Itoa(i) + "]:" + err.Error()
                return errors.New(errMsg)
            }
        }
    }

    return nil
}
