/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "reflect"
    "strings"
    "testing"

    "gopkg.in/yaml.v2"

    factory "../factory"
    // flog "../flog"
    . "../model"
)

const configFile string = "../test/conf/test_nssf_config.yaml"

func initConfigFactory(f string) {
    content, _ := ioutil.ReadFile(f)

    yaml.Unmarshal([]byte(content), &factory.NssfConfig)

    factory.NssfConfig.CheckIntegrity()
}

func generateNonRoamingQueryParameter() NsselectionQueryParameter {
    const jsonQuery = `
        {
            "nf-type": "AMF",
            "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
            "slice-info-request-for-registration": {
                "subscribedNssai": [
                    {
                        "subscribedSnssai": {
                            "sst": 1,
                            "sd": "1"
                        },
                        "defaultIndication": true
                    },
                    {
                        "subscribedSnssai": {
                            "sst": 1,
                            "sd": "3"
                        },
                        "defaultIndication": false
                    }
                ],
                "allowedNssaiCurrentAccess": {
                    "allowedSnssaiList": [
                        {
                            "allowedSnssai": {
                                "sst": 1,
                                "sd": "1"
                            },
                            "nsiInformationList": [
                                {
                                    "nrfId": "http://140.113.194.229:8080/nnrf-nfm/v1/nf-instances",
                                    "nsiId": "1"
                                }
                            ]
                        }
                    ],
                    "accessType": "3GPP_ACCESS"
                },
                "requestedNssai": [
                    {
                        "sst": 1,
                        "sd": "1"
                    },
                    {
                        "sst": 1,
                        "sd": "2"
                    },
                    {
                        "sst": 1,
                        "sd": "3"
                    }
                ],
                "defaultConfiguredSnssaiInd": false
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
}

func generateRoamingQueryParameter() NsselectionQueryParameter {
    const jsonQuery = `
        {
            "nf-type": "AMF",
            "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
            "slice-info-request-for-registration": {
                "subscribedNssai": [
                    {
                        "subscribedSnssai": {
                            "sst": 1,
                            "sd": "1"
                        },
                        "defaultIndication": true
                    },
                    {
                        "subscribedSnssai": {
                            "sst": 1,
                            "sd": "2"
                        },
                        "defaultIndication": false
                    }
                ],
                "allowedNssaiCurrentAccess": {
                    "allowedSnssaiList": [
                        {
                            "allowedSnssai": {
                                "sst": 1,
                                "sd": "1"
                            },
                            "nsiInformationList": [
                                {
                                    "nrfId": "http://140.113.194.229:8080/nnrf-nfm/v1/nf-instances",
                                    "nsiId": "1"
                                }
                            ],
                            "mappedHomeSnssai": {
                                "sst": 1,
                                "sd": "1"
                            }
                        }
                    ],
                    "accessType": "3GPP_ACCESS"
                },
                "requestedNssai": [
                    {
                        "sst": 1,
                        "sd": "1"
                    },
                    {
                        "sst": 1,
                        "sd": "2"
                    },
                    {
                        "sst": 1,
                        "sd": "3"
                    }
                ],
                "defaultConfiguredSnssaiInd": false,
                "mappingofNssai": [
                    {
                        "servingSnssai": {
                            "sst": 1,
                            "sd": "1"
                        },
                        "homeSnssai": {
                            "sst": 1,
                            "sd": "1"
                        }
                    },
                    {
                        "servingSnssai":{
                            "sst": 1,
                            "sd": "2"
                        },
                        "homeSnssai": {
                            "sst": 1,
                            "sd": "3"
                        }
                    },
                    {
                        "servingSnssai":{
                            "sst": 1,
                            "sd": "3"
                        },
                        "homeSnssai": {
                            "sst": 1,
                            "sd": "4"
                        }
                    }
                ]
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
}

func setupUnsupportedHomePlmnId(p *NsselectionQueryParameter) {
    *p.HomePlmnId = PlmnId {
        Mcc: "123",
        Mnc: "456",
    }
}

func setupUnsupportedTai(p *NsselectionQueryParameter) {
    *p.Tai = Tai {
        PlmnId: &PlmnId {
            Mcc: "466",
            Mnc: "92",
        },
        Tac: "12345",
    }
}

func setupUnsupportedSnssai(p *NsselectionQueryParameter) {
    p.SliceInfoRequestForRegistration.RequestedNssai = append(p.SliceInfoRequestForRegistration.RequestedNssai, Snssai {
        Sst: 9,
        Sd: "9",
    })
}

func removeRequestedNssai(p *NsselectionQueryParameter) {
    p.SliceInfoRequestForRegistration.RequestedNssai = []Snssai{}
}

func TestTemplate(t *testing.T) {
    t.Skip()

    initConfigFactory(configFile)

    d, _ := yaml.Marshal(&factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestAddAmfInformation(t *testing.T) {
    initConfigFactory(configFile)
    // flog.MuteLog()

    subtests := []struct {
        name string
        authorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
        expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
    }{
        {
            name: "Add Candidate AMF List",
            authorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 2,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8085/nnrf-nfm/v1/nf-instances",
                                        NsiId: "3",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
            },
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 2,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8085/nnrf-nfm/v1/nf-instances",
                                        NsiId: "3",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
                },
            },
        },
        {
            name: "Add Target AMF Set",
            authorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "3",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8084/nnrf-nfm/v1/nf-instances",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
            },
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "3",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8084/nnrf-nfm/v1/nf-instances",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                TargetAmfSet: "2",
                NrfAmfSet: "http://140.113.194.229:8084/nnrf-nfm/v1/nf-instances",
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            a := *subtest.authorizedNetworkSliceInfo
            addAmfInformation(&a)

            if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) != true {
                e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                r, _ := json.Marshal(&a)
                t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }
}

func TestNonRoaming(t *testing.T) {
    initConfigFactory(configFile)
    // flog.MuteLog()

    subtests := []struct {
        name string
        modifyQueryParameter func(*NsselectionQueryParameter)
        expectStatus int
        expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Unsupported TA",
            modifyQueryParameter: setupUnsupportedTai,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                RejectedNssaiInTa: []Snssai {
                    {
                        Sst: 1,
                        Sd: "1",
                    },
                    {
                        Sst: 1,
                        Sd: "2",
                    },
                    {
                        Sst: 1,
                        Sd: "3",
                    },
                },
            },
        },
        {
            name: "Unsupported S-NSSAI",
            modifyQueryParameter: setupUnsupportedSnssai,
            expectStatus: http.StatusForbidden,
            expectProblemDetails: &ProblemDetails {
                Title: UNSUPPORTED_RESOURCE,
                Status: http.StatusForbidden,
                Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
                Cause: "SNSSAI_NOT_SUPPORTED",
            },
        },
        {
            name: "Requested with Two Rejected in Both PLMN and TA and One Allowed",
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
                },
                RejectedNssaiInPlmn: []Snssai {
                    {
                        Sst: 1,
                        Sd: "2",
                    },
                },
                RejectedNssaiInTa: []Snssai {
                    {
                        Sst: 1,
                        Sd: "3",
                    },
                },
            },
        },
        {
            name: "No Requested and One Subscribed Marked as Default",
            modifyQueryParameter: removeRequestedNssai,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
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

            p := generateNonRoamingQueryParameter()

            if subtest.modifyQueryParameter != nil {
                subtest.modifyQueryParameter(&p)
            }

            status = nsselectionForRegistration(p, &a, &d)

            if status != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
            }

            if status == http.StatusOK {
                if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) == false {
                    e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                    r, _ := json.Marshal(&a)
                    t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            } else {
                if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
                    e, _ := json.Marshal(subtest.expectProblemDetails)
                    r, _ := json.Marshal(&d)
                    t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            }
        })
    }
}

func TestRoaming(t *testing.T) {
    initConfigFactory(configFile)
    // flog.MuteLog()

    subtests := []struct {
        name string
        modifyQueryParameter func(*NsselectionQueryParameter)
        expectStatus int
        expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Unsupported HPLMN",
            modifyQueryParameter: setupUnsupportedHomePlmnId,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                RejectedNssaiInPlmn: []Snssai {
                    {
                        Sst: 1,
                        Sd: "1",
                    },
                    {
                        Sst: 1,
                        Sd: "2",
                    },
                    {
                        Sst: 1,
                        Sd: "3",
                    },
                },
            },
        },
        {
            name: "Unsupported TA",
            modifyQueryParameter: setupUnsupportedTai,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                RejectedNssaiInTa: []Snssai {
                    {
                        Sst: 1,
                        Sd: "1",
                    },
                    {
                        Sst: 1,
                        Sd: "2",
                    },
                    {
                        Sst: 1,
                        Sd: "3",
                    },
                },
            },
        },
        {
            name: "Unsupported S-NSSAI",
            modifyQueryParameter: setupUnsupportedSnssai,
            expectStatus: http.StatusForbidden,
            expectProblemDetails: &ProblemDetails {
                Title: UNSUPPORTED_RESOURCE,
                Status: http.StatusForbidden,
                Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
                Cause: "SNSSAI_NOT_SUPPORTED",
            },
        },
        {
            name: "Requested with Two Rejected in Both PLMN and TA and One Allowed",
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                                MappedHomeSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
                },
                RejectedNssaiInPlmn: []Snssai {
                    {
                        Sst: 1,
                        Sd: "2",
                    },
                },
                RejectedNssaiInTa: []Snssai {
                    {
                        Sst: 1,
                        Sd: "3",
                    },
                },
            },
        },
        {
            name: "No Requested and One Subscribed Marked as Default",
            modifyQueryParameter: removeRequestedNssai,
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo {
                AllowedNssaiList: []AllowedNssai {
                    {
                        AllowedSnssaiList: []AllowedSnssai {
                            {
                                AllowedSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                                NsiInformationList: []NsiInformation {
                                    {
                                        NrfId: "http://140.113.194.229:8081/nnrf-nfm/v1/nf-instances",
                                        NsiId: "1",
                                    },
                                },
                                MappedHomeSnssai: &Snssai {
                                    Sst: 1,
                                    Sd: "1",
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
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

            p := generateRoamingQueryParameter()

            if subtest.modifyQueryParameter != nil {
                subtest.modifyQueryParameter(&p)
            }

            status = nsselectionForRegistration(p, &a, &d)

            if status != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
            }

            if status == http.StatusOK {
                if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) == false {
                    e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                    r, _ := json.Marshal(&a)
                    t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            } else {
                if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
                    e, _ := json.Marshal(subtest.expectProblemDetails)
                    r, _ := json.Marshal(&d)
                    t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
                }
            }
        })
    }
}
