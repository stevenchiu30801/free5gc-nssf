/*
 * NSSF Configuration Factory
 */

package factory

import (
    "fmt"
)

type Configuration struct {

    SupportedNssaiInPlmnList []SupportedNssaiInPlmn `yaml:"supportedNssaiInPlmnList"`

    NsiList []NsiConfig `yaml:"nsiList,omitempty"`

    AmfSetList []AmfSetConfig `yaml:"amfSetList"`

    AmfList []AmfConfig `yaml:"amfList"`

    TaList []TaConfig `yaml:"taList"`

    MappingListFromPlmn []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}

func (c *Configuration) checkIntegrity() error {
    if c.SupportedNssaiInPlmnList == nil || len(c.SupportedNssaiInPlmnList) == 0 {
        return fmt.Errorf("`supportedNssaiInPlmnList` should not be empty")
    } else {
        for i, supportedNssaiInPlmn := range c.SupportedNssaiInPlmnList {
            err := supportedNssaiInPlmn.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`supportedNssaiInPlmnList`[%d]:%s", i, err.Error())
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
        return fmt.Errorf("`amfSetList` should not be empty")
    } else {
        for i, amfSetConfig := range c.AmfSetList {
            err := amfSetConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`amfSetList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.AmfList == nil || len(c.AmfList) == 0 {
        return fmt.Errorf("`amfList` should not be empty")
    } else {
        for i, amfConfig := range c.AmfList {
            err := amfConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`amfList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.TaList == nil || len(c.TaList) == 0 {
        return fmt.Errorf("`taList` should not be empty")
    } else {
        for i, taConfig := range c.TaList {
            err := taConfig.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`taList`[%d]:%s", i, err.Error())
            }
        }
    }

    if c.MappingListFromPlmn == nil || len(c.MappingListFromPlmn) == 0 {
        return fmt.Errorf("`mappingListFromPlmn` should not be empty")
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
