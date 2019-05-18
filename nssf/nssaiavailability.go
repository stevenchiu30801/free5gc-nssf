/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf

import (
    "net/http"
    "reflect"

    factory "../factory"
    . "../model"
)

func nssaiavailabilityPut(nfId string,
                          n NssaiAvailabilityInfo,
                          a *AuthorizedNssaiAvailabilityInfo,
                          d *ProblemDetails) (status int) {
    hitAmf := false
    // Find AMF configuration of given NfId
    // If found, then update the SupportedNssaiAvailabilityData
    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            for _, s := range n.SupportedNssaiAvailabilityData {
                hitTai := false
                for j, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
                    if reflect.DeepEqual(supportedNssaiAvailabilityData.Tai, s.Tai) == true {
                        // Replace SupportedNssaiAvailabilityData if record of the TAI exists
                        factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j] = s

                        hitTai = true
                        break
                    }
                }
                if hitTai == false {
                    // Create new record of the TAI when no corresponding SupportedNssaiAvailability exists
                    factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData =
                        append(factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData, s)
                }
            }
            hitAmf = true
            break
        }
    }

    // If no AMF record is found, create a new one
    if hitAmf == false {
        var amfConfig factory.AmfConfig
        amfConfig.NfId = nfId
        amfConfig.SupportedNssaiAvailabilityData = n.SupportedNssaiAvailabilityData
        factory.NssfConfig.Configuration.AmfList = append(factory.NssfConfig.Configuration.AmfList,
                                                          amfConfig)
    }

    return http.StatusOK
}
