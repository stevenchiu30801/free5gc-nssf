/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    . "../model"
)

type Service string

// List of NSSF service type
const (
    NSSF_NSSELECTION = "Nnssf_NSSelection"
    NSSF_NSSAIAVAILABILITY = "Nnssf_NSSAIAvailability"
)

type AmfSet struct {
    NfId string `yaml:"nfId"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type TaSet struct {
    Tac string `yaml:"tac"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type MappingSet struct {
    OperatorName string `yaml:"operatorName"`
    HomePlmnId *PlmnId `yaml:"homePlmnId"`
    MappingOfSnssai []MappingOfSnssai `yaml:"mappingOfSnssai"`
}

type Info struct {
    Service *Service `yaml:"service"`
    Version string `yaml:"version,omitempty"`
    Title string `yaml:"title,omitempty"`
    Description string `yaml:"description,omitempty"`
    Url string `yaml:"url"`
}

type Configuration struct {
    SupportedSnssaiInPlmn []Snssai `yaml:"supportedSnssaiInPlmn"`
    AmfSet []AmfSet `yaml:"amfSet"`
    TaSet []TaSet `yaml:"taSet"`
    MappingSet []MappingSet `yaml:"mappingSet"`
}

type Config struct {
    Info *Info `yaml:"info"`
    Configuration *Configuration `yaml:"configuration"`
}
