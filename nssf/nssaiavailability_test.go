/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "encoding/json"
    "net/http"
    "reflect"
    "strings"
    "testing"

    "gopkg.in/yaml.v2"

    factory "../factory"
    flog "../flog"
    . "../model"
    test "../test"
)

var testingNssaiavailability = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
    NfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
    GeneratePutRequestBody: func() NssaiAvailabilityInfo {
        const jsonRequest = `
            {
                "supportedNssaiAvailabilityData": [
                    {
                        "tai": {
                            "plmnId": {
                                "mcc": "466",
                                "mnc": "92"
                            },
                            "tac": "33456"
                        },
                        "supportedSnssaiList": [
                            {
                                "sst": 1
                            },
                            {
                                "sst": 1,
                                "sd": "1"
                            },
                            {
                                "sst": 1,
                                "sd": "2"
                            }
                        ]
                    },
                    {
                        "tai": {
                            "plmnId": {
                                "mcc": "466",
                                "mnc": "92"
                            },
                            "tac": "33458"
                        },
                        "supportedSnssaiList": [
                            {
                                "sst": 1
                            },
                            {
                                "sst": 1,
                                "sd": "1"
                            },
                            {
                                "sst": 1,
                                "sd": "3"
                            }
                        ]
                    }
                ],
                "supportedFeatures": ""
            }
        `

        var n NssaiAvailabilityInfo
        json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n)

        return n
    },
}

func generateAddRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "add",
                "path": "/0/supportedSnssaiList/-",
                "value": {
                    "sst": 1,
                    "sd": "1"
                }
            }
        ]
    `

    return []byte(jsonRequest)
}

func generateCopyRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "copy",
                "path": "/1/supportedSnssaiList/2",
                "from": "/0/supportedSnssaiList/2"
            }
        ]
    `

    return []byte(jsonRequest)
}

func generateMoveRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "move",
                "path": "/0/supportedSnssaiList/1",
                "from": "/0/supportedSnssaiList/3"
            }
        ]
    `

    return []byte(jsonRequest)
}

func generateRemoveRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "remove",
                "path": "/0/supportedSnssaiList/2"
            }
        ]
    `

    return []byte(jsonRequest)
}

func generateReplaceRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "replace",
                "path": "/1/supportedSnssaiList/0",
                "value": {
                    "sst": 1,
                    "sd": ""
                }
            }
        ]
    `

    return []byte(jsonRequest)
}

func generateTestRequest() []byte {
    const jsonRequest = `
        [
            {
                "op": "test",
                "path": "/1/supportedSnssaiList/2",
                "value": {
                    "sst": 2,
                    "sd": ""
                }
            }
        ]
    `

    return []byte(jsonRequest)
}

func TestNssaiavailabilityTemplate(t *testing.T) {
    t.Skip()

    // Tests may have different configuration files
    factory.InitConfigFactory(testingNssaiavailability.ConfigFile)

    d, _ := yaml.Marshal(*factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestNssaiavailabilityPatch(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailability.ConfigFile)
    if testingNssaiavailability.MuteLogInd == true {
        flog.Nsselection.MuteLog()
    }

    subtests := []struct {
        name string
        generateRequestBody func() []byte
        expectStatus int
        expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Add",
            generateRequestBody: generateAddRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                        },
                    },
                },
            },
        },
        {
            name: "Copy",
            generateRequestBody: generateCopyRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                },
            },
        },
        {
            name: "Move",
            generateRequestBody: generateMoveRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                },
            },
        },
        {
            name: "Remove",
            generateRequestBody: generateRemoveRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                },
            },
        },
        {
            name: "Replace",
            generateRequestBody: generateReplaceRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                },
            },
        },
        {
            name: "Test",
            generateRequestBody: generateTestRequest,
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33457",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                            {
                                Sst: 2,
                            },
                        },
                    },
                },
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                p []byte
                status int
                a AuthorizedNssaiAvailabilityInfo
                d ProblemDetails
            )

            if subtest.generateRequestBody != nil {
                p = subtest.generateRequestBody()
            }

            status = nssaiavailabilityPatch(testingNssaiavailability.NfId, p, &a, &d)

            if status == http.StatusOK {
                if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
                    e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
                    r, _ := json.Marshal(a)
                    t.Errorf("Incorrect authorized NSSAI availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            } else {
                if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
                    e, _ := json.Marshal(*subtest.expectProblemDetails)
                    r, _ := json.Marshal(d)
                    t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            }
        })
    }
}

func TestNssaiavailabilityPut(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailability.ConfigFile)
    if testingNssaiavailability.MuteLogInd == true {
        flog.Nsselection.MuteLog()
    }

    subtests := []struct {
        name string
        modifyRequestBody func(*NssaiAvailabilityInfo)
        expectStatus int
        expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Create and Replace",
            expectStatus: http.StatusOK,
            expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo {
                AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData {
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33456",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "2",
                            },
                        },
                    },
                    {
                        Tai: &Tai {
                            PlmnId: &PlmnId {
                                Mcc: "466",
                                Mnc: "92",
                            },
                            Tac: "33458",
                        },
                        SupportedSnssaiList: []Snssai {
                            {
                                Sst: 1,
                            },
                            {
                                Sst: 1,
                                Sd: "1",
                            },
                            {
                                Sst: 1,
                                Sd: "3",
                            },
                        },
                        RestrictedSnssaiList: []RestrictedSnssai {
                            {
                                HomePlmnId: &PlmnId {
                                    Mcc: "310",
                                    Mnc: "560",
                                },
                                SNssaiList: []Snssai {
                                    {
                                        Sst: 1,
                                        Sd: "3",
                                    },
                                },
                            },
                        },
                    },
                },
                SupportedFeatures: "",
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                status int
                a AuthorizedNssaiAvailabilityInfo
                d ProblemDetails
            )

            n := testingNssaiavailability.GeneratePutRequestBody()

            if subtest.modifyRequestBody != nil {
                subtest.modifyRequestBody(&n)
            }

            status = nssaiavailabilityPut(testingNssaiavailability.NfId, n, &a, &d)

            if status == http.StatusOK {
                if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
                    e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
                    r, _ := json.Marshal(a)
                    t.Errorf("Incorrect authorized NSSAI availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            } else {
                if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
                    e, _ := json.Marshal(*subtest.expectProblemDetails)
                    r, _ := json.Marshal(d)
                    t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            }
        })
    }
}
