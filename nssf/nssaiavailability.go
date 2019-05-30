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

// Add `SupportedSnssaiList` to configuration for NSSAIAvailability PATCH
func patchAddSupportedSnssaiList(nfId string, tai Tai, supportedSnssaiList []Snssai) {
    hitAmf := false
    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            hitTai := false
            for j, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
                if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) == true {
                    for _, snssai := range supportedSnssaiList {
                        if checkSnssaiInNssai(snssai, supportedNssaiAvailabilityData.SupportedSnssaiList) == false {
                            factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j].SupportedSnssaiList =
                                append(factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j].SupportedSnssaiList,
                                       snssai)
                        }
                    }
                    hitTai = true
                    break
                }
            }
            if hitTai == false {
                var s SupportedNssaiAvailabilityData
                s.Tai = new(Tai)
                *s.Tai = tai
                s.SupportedSnssaiList = supportedSnssaiList
                factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData =
                    append(factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData, s)
            }
            hitAmf = true
            break
        }
    }
    if hitAmf == false {
        var a factory.AmfConfig
        a.NfId = nfId

        var s SupportedNssaiAvailabilityData
        s.Tai = new(Tai)
        *s.Tai = tai
        s.SupportedSnssaiList = supportedSnssaiList

        a.SupportedNssaiAvailabilityData = append(a.SupportedNssaiAvailabilityData, s)

        factory.NssfConfig.Configuration.AmfList = append(factory.NssfConfig.Configuration.AmfList, a)
    }
}

// NSSAIAvailability PATCH method
func nssaiavailabilityPatch(nfId string,
                            p PatchDocument,
                            a *AuthorizedNssaiAvailabilityInfo,
                            d *ProblemDetails) (status int) {
    for i, patchItem := range p {
        var (
            taiInPath *Tai
            // taiInFrom *Tai
            snssaiInPath *Snssai
            // snssaiInFrom *Snssai
            supportedSnssaiList []Snssai
        )

        // Parse `path`
        s := strings.Split(patchItem.Path, "/")
        if len(s) == 0 {
            problemDetail := fmt.Sprintf("[Request Body] [%d]:`path` is invalid")
            *d = ProblemDetails {
                Title: INVALID_REQUEST,
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
        switch s[1] {
        case "":
            // Similiar action as NSSAIAvailability PUT method
        case "supportedNssaiAvailabilityData":
            // PATCH for the whole `SupportedNssaiAvailabilityData` is not supported
            // Since NSSF consumer could simply use NSSAIAvailability PUT method to update
            if len(s) > 2 {
                taiInPath = new(Tai)
                err := json.NewDecoder(strings.NewReader(s[2])).Decode(&taiInPath)
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

                if len(s) > 3 {
                    snssaiInPath = new(Snssai)
                    err = json.NewDecoder(strings.NewReader(s[3])).Decode(&snssaiInPath)
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
                }
            }
        default:
        }

        // Parse `Value`
        if patchItem.Value != nil {
            for key := range *patchItem.Value {
                if key == "supportedSnssaiList" {
                    snssaiListVal := reflect.ValueOf((*patchItem.Value)[key])
                    if snssaiListVal.Kind() != reflect.Slice {
                        problemDetail := fmt.Sprintf("[Request Body] [%d]:`value`:`supportedSnssaiList` should be a valid array")
                        *d = ProblemDetails {
                            Title: INVALID_REQUEST,
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

                    for j := 0; j < snssaiListVal.Len(); j++ {
                        // flog.Nssaiavailability.Infof("%v", snssaiListVal.Index(j).Interface())

                        // Convert map interface to json string then to struct
                        e, _ := json.Marshal(snssaiListVal.Index(j).Interface())

                        var snssai Snssai
                        err := json.NewDecoder(strings.NewReader(string(e))).Decode(&snssai)
                        if err != nil {
                            problemDetail := fmt.Sprintf("[Request Body] [%d]:`value`:supportedSnssaiList`[%d] %s", i, j, err.Error())
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

                        if checkSupportedSnssaiInPlmn(snssai) == false {
                            *d = ProblemDetails {
                                Title: UNSUPPORTED_RESOURCE,
                                Status: http.StatusForbidden,
                                Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
                                Cause: "SNSSAI_NOT_SUPPORTED",
                            }

                            status = http.StatusForbidden
                            return
                        }

                        supportedSnssaiList = append(supportedSnssaiList, snssai)
                    }
                }
            }
        }

        // `From` shall be present if the patch operation is "move" or "copy"
        // `Value` shall be present if the patch operation is "add", "replace" or "test"
        // These are verified in integrity check
        switch *patchItem.Op {
        case ADDPatchOperation:
            if taiInPath != nil && len(supportedSnssaiList) != 0 {
                patchAddSupportedSnssaiList(nfId, *taiInPath, supportedSnssaiList)

                authorizedNssaiAvailabilityData, err := getAuthorizedNssaiAvailabilityDataFromConfig(nfId, *taiInPath)
                if err == nil {
                    a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
                } else {
                    flog.Nssaiavailability.Warnf(err.Error())
                }
            } else {
                *d = ProblemDetails {
                    Title: INVALID_REQUEST,
                    Status: http.StatusBadRequest,
                    Detail: "Both TAI in `path` and `value`:`supportedSnssaiList` should be provided with `op`:'add' operation",
                }

                status = http.StatusBadRequest
                return
            }
        case COPYPatchOperation:
        case MOVEPatchOperation:
        case REMOVEPatchOperation:
        case REPLACEPatchOperation:
        case TESTPatchOperation:
        }
    }

    return http.StatusOK
}

// NSSAIAvailability PUT method
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
        authorizedNssaiAvailabilityData, err := getAuthorizedNssaiAvailabilityDataFromConfig(nfId, *s.Tai)
        if err == nil {
            a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
        } else {
            flog.Nssaiavailability.Warnf(err.Error())
        }
    }

    return http.StatusOK
}
