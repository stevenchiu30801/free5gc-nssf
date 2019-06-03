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

// Parse `path` in `PatchItem`
// Pass pointer and put value of elements in path if provided
func parsePathInPatchItem(path string) (Tai, Snssai, int, error) {
    var (
        tai Tai
        snssai Snssai
    )

    if string(path[0]) == "/" {
        path = path[1:]
    }

    s := strings.Split(path, "/")
    switch s[0] {
    case "":
        // '/'
    case "supportedNssaiAvailabilityData":
        if len(s) > 1 && s[1] != "" {
            // '/supportedNssaiAvailabilityData/{TAI}'
            err := json.NewDecoder(strings.NewReader(s[1])).Decode(&tai)
            if err != nil {
                return tai, snssai, 0, err
            }

            if len(s) > 2 && s[2] != "" {
                // '/supportedNssaiAvailabilityData/{Tai}/{Snssai}'
                err = json.NewDecoder(strings.NewReader(s[2])).Decode(&snssai)
                if err != nil {
                    return tai, snssai, 0, err
                }
                return tai, snssai, 3, nil
            } else {
                return tai, snssai, 2, nil
            }
        } else {
            return tai, snssai, 1, nil
        }
    default:
    }
    return tai, snssai, 0, nil
}

// Parse `value` in `PatchItem`
func parseValueInPatchItem(value map[string]interface{}) ([]Snssai, error) {
    var supportedSnssaiList []Snssai

    for key := range value {
        switch key {
        case "supportedSnssaiList":
            snssaiListVal := reflect.ValueOf(value[key])
            if snssaiListVal.Kind() != reflect.Slice {
                err := fmt.Errorf("`supportedSnssaiList` should be a valid array")
                return supportedSnssaiList, err
            }

            for i := 0; i < snssaiListVal.Len(); i++ {
                // flog.Nssaiavailability.Infof("%v", snssaiListVal.Index(j).Interface())

                // Convert map interface to json string then to struct
                e, _ := json.Marshal(snssaiListVal.Index(i).Interface())

                var snssai Snssai
                err := json.NewDecoder(strings.NewReader(string(e))).Decode(&snssai)
                if err != nil {
                    err = fmt.Errorf("`supportedSnssaiList`[%d] %s", i, err.Error())
                    return supportedSnssaiList, err
                }

                supportedSnssaiList = append(supportedSnssaiList, snssai)
            }
        default:
        }
    }

    return supportedSnssaiList, nil
}

// Add `SupportedSnssaiList` of the given NF ID and TAI to configuration for NSSAIAvailability PATCH
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
                // No supported S-NSSAI list of the given TAI
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
        // No AMF configuration of the given NF ID
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

// Copy `SupportedNssaiAvailabilityData` from one TAI to another for NSSAIAvailability PATCH
func patchCopySupportedNssaiAvailabilityData(nfId string, toTai Tai, fromTai Tai) {
    var (
        toIndex int = -1
        fromIndex int = -1
    )

    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            for j, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
                if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, toTai) == true {
                    toIndex = j
                    if fromIndex != -1 {
                        break
                    }
                }
                if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, fromTai) == true {
                    fromIndex = j
                    if toIndex != -1 {
                        break
                    }
                }
            }

            if fromIndex == -1 {
                e, _ := json.Marshal(&fromTai)
                flog.Nssaiavailability.Warnf("Provided TAI %s in `from` in request body of PATCH request is not found in configuration", e)
                return
            }

            if toIndex == -1 {
                // No existing TAI in `path`, and therefore create a new one
                var s SupportedNssaiAvailabilityData
                s.Tai = new(Tai)
                *s.Tai = toTai
                s.SupportedSnssaiList =
                    factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[fromIndex].SupportedSnssaiList
                factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData = append(
                    factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData, s)
            } else {
                // Replace existing S-NSSAI list with copied list
                factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[toIndex].SupportedSnssaiList =
                    factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[fromIndex].SupportedSnssaiList
            }
            return
        }
    }
}

// Move `SupportedNssaiAvailabilityData` from one TAI to another for NSSAIAvailability PATCH
func patchMoveSupportedNssaiAvailabilityData(nfId string, toTai Tai, fromTai Tai) {
    patchCopySupportedNssaiAvailabilityData(nfId, toTai, fromTai)

    // Delete supported NSSAI availability data of TAI in `from`
    patchRemoveSupportedNssaiAvailabilityData(nfId, fromTai)
}

// Remove `SupportedNssaiAvailabilityData` of the given TAI for NSSAIAvailability PATCH
func patchRemoveSupportedNssaiAvailabilityData(nfId string, tai Tai) {
    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            for j, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
                if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) == true {
                    factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData = append(
                        factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[:j],
                        factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j + 1:]...)
                }
            }
            break
        }
    }
    return
}

// Remove S-NSSAI in `SupportedSnssaiList` of the given TAI for NSSAIAvailability PATCH
func patchRemoveSupportedSnssai(nfId string, tai Tai, snssai Snssai) {
    for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            for j, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
                if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) == true {
                    for k, supportedSnssai := range supportedNssaiAvailabilityData.SupportedSnssaiList {
                        if supportedSnssai == snssai {
                            factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j].SupportedSnssaiList = append(
                                factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j].SupportedSnssaiList[:k],
                                factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData[j].SupportedSnssaiList[k + 1:]...)
                        }
                    }
                }
            }
            break
        }
    }
    return
}

// NSSAIAvailability PATCH method
func nssaiavailabilityPatch(nfId string,
                            p PatchDocument,
                            a *AuthorizedNssaiAvailabilityInfo,
                            d *ProblemDetails) (status int) {
    for i, patchItem := range p {
        var (
            taiInPath Tai
            taiInFrom Tai
            snssaiInPath Snssai
            snssaiInFrom Snssai
            depthInPath int
            depthInFrom int
            supportedSnssaiList []Snssai
        )

        // Parse `path`
        taiInPath, snssaiInPath, depthInPath, err := parsePathInPatchItem(patchItem.Path)
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

        // PATCH for the whole `SupportedNssaiAvailabilityData` is not supported
        // Since NSSF consumer could simply use NSSAIAvailability PUT/DELETE method to update
        // Therefore TAI must be provided in `path` on NSSAIAvailability PATCH method
        if depthInPath < 2 {
            problemDetail := fmt.Sprintf("[Request Body] [%d]:`path` TAI should be provided in path", i)
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

        // Parse `From`
        if patchItem.From != "" {
            taiInFrom, snssaiInFrom, depthInFrom, err = parsePathInPatchItem(patchItem.From)
            if err != nil {
                problemDetail := fmt.Sprintf("[Request Body] [%d]:`from` %s", i, err.Error())
                *d = ProblemDetails {
                    Title: MALFORMED_REQUEST,
                    Status: http.StatusBadRequest,
                    Detail: problemDetail,
                    InvalidParams: []InvalidParam {
                        {
                            Param: "from",
                            Reason: problemDetail,
                        },
                    },
                }

                status = http.StatusBadRequest
                return
            }

            // Structure of `From` should match that of `Path` if `From` is provided
            if depthInPath != depthInFrom {
                problemDetail := fmt.Sprintf("[Request Body] [%d]:`path` and `from` should have same path structure", i)
                *d = ProblemDetails {
                    Title: INVALID_REQUEST,
                    Status: http.StatusBadRequest,
                    Detail: problemDetail,
                    InvalidParams: []InvalidParam {
                        {
                            Param: "path",
                            Reason: problemDetail,
                        },
                        {
                            Param: "from",
                            Reason: problemDetail,
                        },
                    },
                }

                status = http.StatusBadRequest
                return
            }
        }

        // Parse `Value`
        if patchItem.Value != nil {
            supportedSnssaiList, err = parseValueInPatchItem(*patchItem.Value)
            if err != nil {
                problemDetail := fmt.Sprintf("[Request Body] [%d]:`value`:%s", i, err.Error())
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

            // Check if all S-NSSAIs is valid in the PLMN
            if checkSupportedNssaiInPlmn(supportedSnssaiList, *taiInPath.PlmnId) == false {
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

        flog.Nssaiavailability.Warnf("Delete this log after %v and %v are used", snssaiInPath, snssaiInFrom)

        // `From` shall be present if the patch operation is "move" or "copy"
        // `Value` shall be present if the patch operation is "add", "replace" or "test"
        // These are verified in integrity check
        switch *patchItem.Op {
        case ADDPatchOperation:
            if len(supportedSnssaiList) == 0 {
                problemDetail := "[Request Body] `value`:`supportedSnssaiList` should not be empty with `op`:'add' operation"
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

            if depthInPath == 2 {
                // Add supported S-NSSAI list
                patchAddSupportedSnssaiList(nfId, taiInPath, supportedSnssaiList)
            }
        case COPYPatchOperation:
            if depthInPath == 2 {
                // Copy supported S-NSSAI list between TAIs
                patchCopySupportedNssaiAvailabilityData(nfId, taiInPath, taiInFrom)
            }
        case MOVEPatchOperation:
            if depthInPath == 2 {
                // Move supported S-NSSAI list between TAIs
                patchMoveSupportedNssaiAvailabilityData(nfId, taiInPath, taiInFrom)
            }
        case REMOVEPatchOperation:
            if depthInPath == 2 {
                // Remove supported S-NSSAI list
                patchRemoveSupportedNssaiAvailabilityData(nfId, taiInPath)
            } else if depthInPath == 3 {
                // Remove specified supported S-NSSAI
                patchRemoveSupportedSnssai(nfId, taiInPath, snssaiInPath)
            }
        case REPLACEPatchOperation:
        case TESTPatchOperation:
        }

        authorizedNssaiAvailabilityData, err := getAuthorizedNssaiAvailabilityDataFromConfig(nfId, taiInPath)
        if err == nil {
            a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
        } else {
            flog.Nssaiavailability.Warnf(err.Error())
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
