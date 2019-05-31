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

// Parse path in `PatchItem`
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
        if len(s) > 1 {
            // '/supportedNssaiAvailabilityData/{TAI}'
            err := json.NewDecoder(strings.NewReader(s[1])).Decode(&tai)
            if err != nil {
                return tai, snssai, 0, err
            }

            if len(s) > 2 {
                // '/supportedNssaiAvailabilityData/{Tai}/{Snssai}'
                err = json.NewDecoder(strings.NewReader(s[2])).Decode(&snssai)
                if err != nil {
                    return tai, snssai, 0, err
                }
                return tai, snssai, 3, nil
            }
            return tai, snssai, 2, nil
        }
    default:
    }
    return tai, snssai, 0, nil
}

func parseValueInPatchItem(value map[string]interface{}) ([]Snssai, error) {
    var supportedSnssaiList []Snssai

    for key := range value {
        if key == "supportedSnssaiList" {
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
        }
    }

    return supportedSnssaiList, nil
}

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
            taiInFrom *Tai
            snssaiInPath *Snssai
            snssaiInFrom *Snssai
            supportedSnssaiList []Snssai
        )

        // Parse `path`
        tai, snssai, depthInPath, err := parsePathInPatchItem(patchItem.Path)
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

        if depthInPath == 2 {
            taiInPath = new(Tai)
            *taiInPath = tai
        }
        if depthInPath == 3 {
            taiInPath = new(Tai)
            *taiInPath = tai
            snssaiInPath = new(Snssai)
            *snssaiInPath = snssai
        }

        // PATCH for the whole `SupportedNssaiAvailabilityData` is not supported
        // Since NSSF consumer could simply use NSSAIAvailability PUT method to update
        // Therefore TAI must be provided in `path` on NSSAIAvailability PATCH method
        if taiInPath == nil {
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
            tai, snssai, depthInFrom, err := parsePathInPatchItem(patchItem.Path)
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

            if depthInFrom == 2 {
                taiInFrom = new(Tai)
                *taiInFrom = tai
            }
            if depthInFrom == 3 {
                taiInFrom = new(Tai)
                *taiInFrom = tai
                snssaiInFrom = new(Snssai)
                *snssaiInFrom = snssai
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
            if checkSupportedNssaiInPlmn(supportedSnssaiList) == false {
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

        // `From` shall be present if the patch operation is "move" or "copy"
        // `Value` shall be present if the patch operation is "add", "replace" or "test"
        // These are verified in integrity check
        switch *patchItem.Op {
        case ADDPatchOperation:
            if len(supportedSnssaiList) != 0 {
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
                    Detail: "[Request Body] `value`:`supportedSnssaiList` should not be empty with `op`:'add' operation",
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
