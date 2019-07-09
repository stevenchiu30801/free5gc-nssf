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

type NsiConfig struct {

    Snssai *Snssai `yaml:"snssai"`

    NsiInformationList []NsiInformation `yaml:"nsiInformationList"`
}

func (n *NsiConfig) checkIntegrity() error {
    if n.Snssai == nil {
        return fmt.Errorf("`snssai` should not be empty")
    } else {
        err := n.Snssai.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`snssai`:%s", err.Error())
        }
    }

    if n.NsiInformationList == nil || len(n.NsiInformationList) == 0 {
        return fmt.Errorf("`nsiInformation` should not be empty")
    } else {
        for i, nsiInformation := range n.NsiInformationList {
            err := nsiInformation.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`nsiInformation`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
