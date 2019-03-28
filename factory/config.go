/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

// import (
//     . "../model"
// )

type Service string

// List of NSSF service type
const (
    NSSF_NSSELECTION = "Nnssf_NSSelection"
    NSSF_NSSAIAVAILABILITY = "Nnssf_NSSAIAvailability"
)

type PlmnId struct {
    Mcc string `yaml:"mcc"`
    Mnc string `taml:"mnc"`
}

type Snssai struct {
    Sst int32 `yaml:"sst"`
    Sd string `yaml:"sd,omitempty"`
}

type Tai struct {
    PlmnId *PlmnId `yaml:"plmnId"`
    Tac string `yaml:"tac"`
}

type AmfSet struct {
    NfId string `yaml:"nfId"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type TaSet struct {
    Tai *Tai `yaml:"tai"`
    SupportedSnssai []Snssai `yaml:"supportedSnssai"`
}

type Info struct {
    Service *Service `yaml:"service"`
    Version string `yaml:"version,omitempty"`
    Title string `yaml:"title,omitempty"`
    Description string `yaml:"description,omitempty"`
    Url string `yaml:"url"`
}

type Configuration struct {
    AmfSet []AmfSet `yaml:"amfSet"`
    TaSet []TaSet `yaml:"taSet"`
}

type Config struct {
    Info *Info `yaml:"info"`
    Configuration *Configuration `yaml:"configuration"`
}
