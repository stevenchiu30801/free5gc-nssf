/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf

import (
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
    "strings"

    factory "../factory"
    flog "../flog"
    . "../model"
)

func nssaiavailabilityPatch(nfId string,
                            p PatchDocument,
                            a *AuthorizedNssaiAvailabilityInfo,
                            d *ProblemDetails) (status int) {
    for i, patchItem := range p {
        s := strings.Split(patchItem.Path, "/")

        switch s[1] {
        case "supportedNssaiAvailabilityData":
            if s[2] != "" {
                var tai Tai
                err := json.NewDecoder(strings.NewReader(s[2])).Decode(&tai)
                if err != nil {
                    problemDetail := fmt.Sprintf("[Request Body] [%d]:`path` %s", i, err.Error())
                    *d = ProblemDetails {
                        Title: MALFORMED_REQUEST,
                        Status: http.StatusBadRequest,
                        Detail: problemDetail,
                        InvalidParams: []InvalidParam {
                            {
                                Param: "path",
                                Reason: problemDetail,
                            },
                        },
                    }

                    status = http.StatusBadRequest
                    return
                }

                for key := range *patchItem.Value {
                    if key == "supportedSnssaiList" {
                        snssaiMap, _ := (*patchItem.Value)[key]
                        snssaiListIntf := reflect.ValueOf(snssaiMap)
                        if snssaiListIntf.Kind() != reflect.Slice {
                            problemDetail := fmt.Sprintf("[Request Body] [%d]:`value`:`supportedSnssaiList` should be a valid slice")
                            *d = ProblemDetails {
                                Title: MALFORMED_REQUEST,
                                Status: http.StatusBadRequest,
                                Detail: problemDetail,
                                InvalidParams: []InvalidParam {
                                    {
                                        Param: "supportedSnssaiList",
                                        Reason: problemDetail,
                                    },
                                },
                            }

                            status = http.StatusBadRequest
                            return
                        }

                        for j := 0; j < snssaiListIntf.Len(); j++ {
                            flog.Nssaiavailability.Infof("%v", snssaiListIntf.Index(j).Interface())
                        }
                    }


                }
            }
        default:
        }
    }

    return http.StatusOK
}

func nssaiavailabilityPut(nfId string,
                          n NssaiAvailabilityInfo,
                          a *AuthorizedNssaiAvailabilityInfo,
                          d *ProblemDetails) (status int) {
    for _, s := range n.SupportedNssaiAvailabilityData {
        if checkSupportedNssaiInPlmn(s.SupportedSnssaiList) == false {
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

    for _, s := range n.SupportedNssaiAvailabilityData {
        var authorizedNssaiAvailabilityData AuthorizedNssaiAvailabilityData
        authorizedNssaiAvailabilityData.Tai = s.Tai
        authorizedNssaiAvailabilityData.SupportedSnssaiList = s.SupportedSnssaiList

        for _, taConfig := range factory.NssfConfig.Configuration.TaList {
            if reflect.DeepEqual(taConfig.Tai, s.Tai) == true {
                if taConfig.RestrictedSnssaiList != nil && len(taConfig.RestrictedSnssaiList) != 0 {
                    authorizedNssaiAvailabilityData.RestrictedSnssaiList = taConfig.RestrictedSnssaiList
                }
                break
            }
        }

        a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
    }

    return http.StatusOK
}
