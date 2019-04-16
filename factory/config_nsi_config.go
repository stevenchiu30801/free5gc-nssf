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

type NsiConfig struct {

    Snssai *Snssai `yaml:"snssai"`

    NsiInformationList []NsiInformation `yaml:"nsiInformationList"`
}

func (n *NsiConfig) checkIntegrity() error {
    if n.Snssai == nil {
        return errors.New("`snssai` in configuration should not be empty")
    } else {
        err := n.Snssai.CheckIntegrity()
        if err != nil {
            errMsg := "`snssai`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if n.NsiInformationList == nil || len(n.NsiInformationList) == 0 {
        return errors.New("`nsiInformation` in configuration should not be empty")
    } else {
        for i, nsiInformation := range n.NsiInformationList {
            err := nsiInformation.CheckIntegrity()
            if err != nil {
                errMsg := "`nsiInformation`[" + strconv.Itoa(i) + "]:" + err.Error()
                return errors.New(errMsg)
            }
        }
    }

    return nil
}
