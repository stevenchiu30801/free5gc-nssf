/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"

    . "free5gc-nssf/model"
)

type SupportedNssaiInPlmn struct {

    PlmnId *PlmnId `yaml:"plmnId"`

    SupportedSnssaiList []Snssai `yaml:"supportedSnssaiList"`
}

func (s *SupportedNssaiInPlmn) checkIntegrity() error {
    if s.PlmnId == nil {
        return fmt.Errorf("`plmnId` should not be empty")
    } else {
        err := s.PlmnId.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`plmnId`:%s", err.Error())
        }
    }

    if s.SupportedSnssaiList == nil || len(s.SupportedSnssaiList) == 0 {
        return fmt.Errorf("`supportedSnssaiList` should not be empty")
    } else {
        for i, snssai := range s.SupportedSnssaiList {
            err := snssai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedSnssaiList", i, err.Error())
            }
        }
    }

    return nil
}
