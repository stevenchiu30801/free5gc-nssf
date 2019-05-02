/*
 * NSSF NS Selection
 *
 * NSsf Network Slice Selection Service
 */

package test

import (
    "flag"

    . "../model"
)

var (
    ConfigFileFromArgs string
    MuteLogIndFromArgs bool
)

type TestingParameter struct {

    ConfigFile string

    MuteLogInd bool

    GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

    GenerateRoamingQueryParameter func() NsselectionQueryParameter
}

func init() {
    flag.StringVar(&ConfigFileFromArgs, "config-file", "../test/conf/test_nssf_config.yaml", "Configuration file")
    flag.BoolVar(&MuteLogIndFromArgs, "mute-log", false, "Mute log indication")
    flag.Parse()
}
