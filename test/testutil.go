/*
 * NSSF NS Selection
 *
 * NSsf Network Slice Selection Service
 */

package test

import (
    . "../model"
)

type TestingParameter struct {

    ConfigFile string

    MuteLogInd bool

    GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

    GenerateRoamingQueryParameter func() NsselectionQueryParameter
}
