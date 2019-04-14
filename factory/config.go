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
    NSSF_NSSELECTION Service = "Nnssf_NSSelection"
    NSSF_NSSAIAVAILABILITY Service = "Nnssf_NSSAIAvailability"
)

type AmfConfig struct {
    NfId string `yaml:"nfId"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type TaConfig struct {
    Tac string `yaml:"tac"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type MappingFromPlmnConfig struct {
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
    ServingPlmnId *PlmnId `yaml:"servingPlmnId"`
}

type Configuration struct {
    SupportedNssaiInPlmn []Snssai `yaml:"supportedNssaiInPlmn"`
    AmfList []AmfConfig `yaml:"amfList"`
    TaList []TaConfig `yaml:"taList"`
    MappingListFromPlmn []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}

type Config struct {
    Info *Info `yaml:"info"`
    Configuration *Configuration `yaml:"configuration"`
}
