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
    GenerateAddRequestBody: func() NssaiAvailabilityInfo {
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
    GeneratePatchRequestBody: func() PatchDocument {
        const jsonRequest = `
            [
                {
                    "op": "add",
                    "path": "/supportedNssaiAvailabilityData/{\"plmnId\":{\"mcc\":\"466\",\"mnc\":\"92\"},\"tac\":\"33456\"}",
                    "value": {
                        "supportedSnssaiList": [
                            {
                                "sst": 1
                            },
                            {
                                "sst": 1,
                                "sd": "1"
                            }
                        ]
                    }
                }
            ]
        `

        var p PatchDocument
        json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p)

        return p
    },
}

func TestNssaiavailabilityTemplate(t *testing.T) {
    t.Skip()

    // Tests may have different configuration files
    factory.InitConfigFactory(testingNssaiavailability.ConfigFile)

    d, _ := yaml.Marshal(*factory.NssfConfig.Info)
    t.Logf("%s", string(d))
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

            n := testingNssaiavailability.GenerateAddRequestBody()

            if subtest.modifyRequestBody != nil {
                subtest.modifyRequestBody(&n)
            }

            status = nssaiavailabilityPut(testingNssaiavailability.NfId, n, &a, &d)

            if status == http.StatusOK {
                a.Sort()
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
