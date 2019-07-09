/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nsselection

import (
    "encoding/json"
    "net/http"
    "reflect"
    "strings"
    "testing"

    "gopkg.in/yaml.v2"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
    "free5gc-nssf/test"
)

var testingNsselectionForPduSession = test.TestingNsselection {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
    GenerateNonRoamingQueryParameter: func() NsselectionQueryParameter {
        const jsonQuery = `
            {
                "nf-type": "AMF",
                "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
                "slice-info-request-for-pdu-session": {
                    "sNssai": {
                        "sst": 1,
                        "sd": "2"
                    },
                    "roamingIndication": "NON_ROAMING"
                },
                "tai": {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                }
            }
        `

        var p NsselectionQueryParameter
        json.NewDecoder(strings.NewReader(jsonQuery)).Decode(&p)

        return p
    },
    GenerateRoamingQueryParameter: func() NsselectionQueryParameter {
        const jsonQuery = `
            {
                "nf-type": "AMF",
                "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
                "slice-info-request-for-pdu-session": {
                    "sNssai": {
                        "sst": 1,
                        "sd": "2"
                    },
                    "roamingIndication": "LOCAL_BREAKOUT",
                    "homeSnssai": {
                        "sst": 1,
                        "sd": "3"
                    }
                },
                "home-plmn-id": {
                    "mcc": "440",
                    "mnc": "10"
                },
                "tai": {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                }
            }
        `

        var p NsselectionQueryParameter
        json.NewDecoder(strings.NewReader(jsonQuery)).Decode(&p)

        return p
    },
}

func setNonRoaming(p *NsselectionQueryParameter) {
    *p.SliceInfoRequestForPduSession.RoamingIndication = RoamingIndication_NON_ROAMING
}

func setLocalBreakout(p *NsselectionQueryParameter) {
    *p.SliceInfoRequestForPduSession.RoamingIndication = RoamingIndication_LOCAL_BREAKOUT
}

func setHomeRoutedRoaming(p *NsselectionQueryParameter) {
    *p.SliceInfoRequestForPduSession.RoamingIndication = RoamingIndication_HOME_ROUTED_ROAMING
}

func checkAuthorizedNetworkSliceInfo(target AuthorizedNetworkSliceInfo, expectList []AuthorizedNetworkSliceInfo) bool {
    target.Sort()
    for _, expectElement := range expectList {
        if reflect.DeepEqual(target, expectElement) == true {
            return true
        }
    }
    return false
}

func TestNsselectionForPduSessionTemplate(t *testing.T) {
    t.Skip()

    // Tests may have different configuration files
    factory.InitConfigFactory(testingNsselectionForPduSession.ConfigFile)

    d, _ := yaml.Marshal(*factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestNsselectionForPduSessionNonRoaming(t *testing.T) {
    factory.InitConfigFactory(testingNsselectionForPduSession.ConfigFile)
    if testingNsselectionForPduSession.MuteLogInd == true {
        flog.Nsselection.MuteLog()
    }

    subtests := []struct {
        name string
        modifyQueryParameter func(*NsselectionQueryParameter)
        expectStatus int
        expectAuthorizedNetworkSliceInfo []AuthorizedNetworkSliceInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Non Roaming",
            modifyQueryParameter: setNonRoaming,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: []AuthorizedNetworkSliceInfo {
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-1.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-2.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                status int
                a AuthorizedNetworkSliceInfo
                d ProblemDetails
            )

            p := testingNsselectionForPduSession.GenerateNonRoamingQueryParameter()

            if subtest.modifyQueryParameter != nil {
                subtest.modifyQueryParameter(&p)
            }

            status = nsselectionForPduSession(p, &a, &d)

            if status != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
            }

            if status == http.StatusOK {
                if checkAuthorizedNetworkSliceInfo(a, subtest.expectAuthorizedNetworkSliceInfo) == false {
                    e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                    r, _ := json.Marshal(a)
                    t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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

func TestNsselectionForPduSessionRoaming(t *testing.T) {
    factory.InitConfigFactory(testingNsselectionForPduSession.ConfigFile)
    if testingNsselectionForPduSession.MuteLogInd == true {
        flog.Nsselection.MuteLog()
    }

    subtests := []struct {
        name string
        modifyQueryParameter func(*NsselectionQueryParameter)
        expectStatus int
        expectAuthorizedNetworkSliceInfo []AuthorizedNetworkSliceInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Local Breakout",
            modifyQueryParameter: setLocalBreakout,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: []AuthorizedNetworkSliceInfo {
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-1.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-2.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },

            },
        },
        {
            name: "Home Routed Roaming",
            modifyQueryParameter: setHomeRoutedRoaming,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: []AuthorizedNetworkSliceInfo {
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-1.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },
                {
                    NsiInformation: &NsiInformation {
                        NrfId: "http://free5gc-nrf-12-2.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                        NsiId: "12",
                    },
                },

            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                status int
                a AuthorizedNetworkSliceInfo
                d ProblemDetails
            )

            p := testingNsselectionForPduSession.GenerateRoamingQueryParameter()

            if subtest.modifyQueryParameter != nil {
                subtest.modifyQueryParameter(&p)
            }

            status = nsselectionForPduSession(p, &a, &d)

            if status != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
            }

            if status == http.StatusOK {
                if checkAuthorizedNetworkSliceInfo(a, subtest.expectAuthorizedNetworkSliceInfo) == false {
                    e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                    r, _ := json.Marshal(a)
                    t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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
