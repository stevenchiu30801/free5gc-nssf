/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf

import (
    "encoding/json"
    "net/http"

    jsonpatch "github.com/evanphx/json-patch"

    factory "../factory"
    flog "../flog"
    . "../model"
)

// NSSAIAvailability PATCH method
func nssaiavailabilityPatch(nfId string,
                            patchDocument []byte,
                            a *AuthorizedNssaiAvailabilityInfo,
                            d *ProblemDetails) (status int) {
    var amfIdx int
    var original []byte
    for amfIdx, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            original, _ = json.Marshal(factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)
            break
        }
    }

    // TODO: Check if returned HTTP status codes or problem details are proper when errors occur

    patch, err := jsonpatch.DecodePatch(patchDocument)
    if err != nil {
        *d = ProblemDetails {
            Title: MALFORMED_REQUEST,
            Status: http.StatusBadRequest,
            Detail: err.Error(),
        }

        status = http.StatusBadRequest
        return
    }

    modified, err := patch.Apply(original)
    if err != nil {
        *d = ProblemDetails {
            Title: INVALID_REQUEST,
            Status: http.StatusConflict,
            Detail: err.Error(),
        }

        status = http.StatusConflict
        return
    }

    err = json.Unmarshal(modified, &factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)
    if err != nil {
        *d = ProblemDetails {
            Title: INVALID_REQUEST,
            Status: http.StatusBadRequest,
            Detail: err.Error(),
        }

        status = http.StatusBadRequest
        return
    }

    // Return all authorized NSSAI availability information
    a.AuthorizedNssaiAvailabilityData, _ = getAllAuthorizedNssaiAvailabilityDataFromConfig(nfId)

    // TODO: Return authorized NSSAI availability information of updated TAI only

    return http.StatusOK
}

// NSSAIAvailability PUT method
func nssaiavailabilityPut(nfId string,
                          n NssaiAvailabilityInfo,
                          a *AuthorizedNssaiAvailabilityInfo,
                          d *ProblemDetails) (status int) {
    for _, s := range n.SupportedNssaiAvailabilityData {
        if checkSupportedNssaiInPlmn(s.SupportedSnssaiList, *s.Tai.PlmnId) == false {
            *d = ProblemDetails {
                Title: UNSUPPORTED_RESOURCE,
                Status: http.StatusForbidden,
                Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
                Cause: "SNSSAI_NOT_SUPPORTED",
            }

            status = http.StatusForbidden
            return
        }
    }

    // TODO: Currently authorize all the provided S-NSSAIs
    //       Take some issue into consideration e.g. operator policies

    hitAmf := false
    // Find AMF configuration of given NfId
    // If found, then update the SupportedNssaiAvailabilityData
    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData = n.SupportedNssaiAvailabilityData

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

    // Return all authorized NSSAI availability information
    // a.AuthorizedNssaiAvailabilityData, _ = getAllAuthorizedNssaiAvailabilityDataFromConfig(nfId)

    // Return authorized NSSAI availability information of updated TAI only
    for _, s := range n.SupportedNssaiAvailabilityData {
        authorizedNssaiAvailabilityData, err := getAuthorizedNssaiAvailabilityDataFromConfig(nfId, *s.Tai)
        if err == nil {
            a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
        } else {
            flog.Nssaiavailability.Warnf(err.Error())
        }
    }

    return http.StatusOK
}
