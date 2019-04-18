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

type Configuration struct {

    SupportedNssaiInPlmn []Snssai `yaml:"supportedNssaiInPlmn"`

    NsiList []NsiConfig `yaml:"nsiList,omitempty"`

    AmfSetList []AmfSetConfig `yaml:"amfSetList"`

    TaList []TaConfig `yaml:"taList"`

    MappingListFromPlmn []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}

func (c *Configuration) checkIntegrity() error {
    if c.SupportedNssaiInPlmn == nil || len(c.SupportedNssaiInPlmn) == 0 {
        return fmt.Errorf("`supportedNssaiInPlmn` in configuration should not be empty")
    } else {
        for i, supportedSnssaiInPlmn := range c.SupportedNssaiInPlmn {
            err := supportedSnssaiInPlmn.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedNssaiInPlmn`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.NsiList != nil && len(c.NsiList) != 0 {
        for i, nsiConfig := range c.NsiList {
            err := nsiConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`nsiList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.AmfSetList == nil || len(c.AmfSetList) == 0 {
        return fmt.Errorf("`amfSetList` in configuration should not be empty")
    } else {
        for i, amfSetConfig := range c.AmfSetList {
            err := amfSetConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`amfSetList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.TaList == nil || len(c.TaList) == 0 {
        return fmt.Errorf("`taList` in configuration should not be empty")
    } else {
        for i, taConfig := range c.TaList {
            err := taConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`taList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.MappingListFromPlmn == nil || len(c.MappingListFromPlmn) == 0 {
        return fmt.Errorf("`mappingListFromPlmn` in configuration should not be empty")
    } else {
        for i, mappingFromPlmnConfig := range c.MappingListFromPlmn {
            err := mappingFromPlmnConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`mappingListFromPlmn`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
