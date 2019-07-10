/*
 * NSSF Testing Utility
 */

package test

import (
    "flag"

    . "free5gc-nssf/model"
)

var (
    ConfigFileFromArgs string
    MuteLogIndFromArgs bool
)

type TestingUtil struct {

    ConfigFile string

    MuteLogInd bool
}

type TestingNsselection struct {

    ConfigFile string

    MuteLogInd bool

    GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

    GenerateRoamingQueryParameter func() NsselectionQueryParameter
}

type TestingNssaiavailability struct {

    ConfigFile string

    MuteLogInd bool

    NfId string

    SubscriptionId string

    NfNssaiAvailabilityUri string
}

func init() {
    flag.StringVar(&ConfigFileFromArgs, "config-file", "../test/conf/test_nssf_config.yaml", "Configuration file")
    flag.BoolVar(&MuteLogIndFromArgs, "mute-log", false, "Mute log indication")
    flag.Parse()
}
