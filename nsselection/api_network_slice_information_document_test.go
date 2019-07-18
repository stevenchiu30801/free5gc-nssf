/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nsselection

import (
    "context"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/antihax/optional"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
    "free5gc-nssf/nsselection/client"
    "free5gc-nssf/test"
)

var testingNsselectionApi = test.TestingNsselection {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
}

func TestNSSelectionGet(t *testing.T) {
    factory.InitConfigFactory(testingNsselectionApi.ConfigFile)
    if testingNsselectionApi.MuteLogInd == true {
        flog.Nsselection.MuteLog()
        flog.Util.MuteLog()
    }

    go func() {
        router := NewRouter()
        err := router.Run(":8080")
        if err != nil {
            t.Errorf(err.Error())
        }
    }()

    configuration := client.NewConfiguration()
    configuration.SetBasePath("http://localhost:8080")
    apiClient := client.NewAPIClient(configuration)

    subtests := []struct {
        name string
        nfType NfType
        nfId string
        localVarOptionals *client.NSSelectionGetParamOpts
        expectStatus int
        expectAuthorizedNetworkSliceInfo []AuthorizedNetworkSliceInfo
    }{
        {
            name: "For Registration",
            nfType: NfType_AMF,
            nfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
            localVarOptionals: &client.NSSelectionGetParamOpts {
                SliceInfoRequestForRegistration: client.ConvertQueryParamIntf(SliceInfoForRegistration {
                    SubscribedNssai: []SubscribedSnssai {
                        {
                            SubscribedSnssai: &Snssai {
                                Sst: 1,
                            },
                            DefaultIndication: false,
                        },
                        {
                            SubscribedSnssai: &Snssai {
                                Sst: 1,
                                Sd: "1",
                            },
                            DefaultIndication: true,
                        },
                    },
                    RequestedNssai: []Snssai {
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
                    MappingOfNssai: []MappingOfSnssai {
                        {
                            ServingSnssai: &Snssai {
                                Sst: 1,
                                Sd: "1",
                            },
                            HomeSnssai: &Snssai {
                                Sst: 1,
                                Sd: "1",
                            },
                        },
                        {
                            ServingSnssai: &Snssai {
                                Sst: 1,
                                Sd: "2",
                            },
                            HomeSnssai: &Snssai {
                                Sst: 1,
                                Sd: "3",
                            },
                        },
                        {
                            ServingSnssai: &Snssai {
                                Sst: 1,
                                Sd: "3",
                            },
                            HomeSnssai: &Snssai {
                                Sst: 1,
                                Sd: "4",
                            },
                        },
                    },
                }),
                SliceInfoRequestForPduSession: optional.EmptyInterface(),
                HomePlmnId: client.ConvertQueryParamIntf(PlmnId {
                    Mcc: "440",
                    Mnc: "10",
                }),
                Tai: client.ConvertQueryParamIntf(Tai {
                    PlmnId: &PlmnId {
                        Mcc: "466",
                        Mnc: "92",
                    },
                    Tac: "33456",
                }),
                SupportedFeatures: optional.EmptyString(),
            },
            expectStatus: http.StatusOK,
            expectAuthorizedNetworkSliceInfo: []AuthorizedNetworkSliceInfo {
                {
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
                                            NrfId: "http://free5gc-nrf-11.nctu.me:8080/nnrf-nfm/v1/nf-instances",
                                            NsiId: "11",
                                        },
                                    },
                                    MappedHomeSnssai: &Snssai {
                                        Sst: 1,
                                        Sd: "1",
                                    },
                                },
                            },
                            AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
                        },
                    },
                    ConfiguredNssai: []ConfiguredSnssai {
                        {
                            ConfiguredSnssai: &Snssai {
                                Sst: 1,
                            },
                        },
                        {
                            ConfiguredSnssai: &Snssai {
                                Sst: 1,
                                Sd: "1",
                            },
                            MappedHomeSnssai: &Snssai {
                                Sst: 1,
                                Sd: "1",
                            },
                        },
                    },
                    CandidateAmfList: []string {
                        "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                        "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
                        "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
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
        },
        {
            name: "For PDU Session",
            nfType: NfType_AMF,
            nfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
            localVarOptionals: &client.NSSelectionGetParamOpts {
                SliceInfoRequestForRegistration: optional.EmptyInterface(),
                SliceInfoRequestForPduSession: client.ConvertQueryParamIntf(SliceInfoForPduSession {
                    SNssai: &Snssai {
                        Sst: 1,
                        Sd: "2",
                    },
                    RoamingIndication: func() RoamingIndication { r := RoamingIndication_HOME_ROUTED_ROAMING; return r }(),
                    HomeSnssai: &Snssai {
                        Sst: 1,
                        Sd: "3",
                    },
                }),
                HomePlmnId: client.ConvertQueryParamIntf(PlmnId {
                    Mcc: "440",
                    Mnc: "10",
                }),
                Tai: client.ConvertQueryParamIntf(Tai {
                    PlmnId: &PlmnId {
                        Mcc: "466",
                        Mnc: "92",
                    },
                    Tac: "33456",
                }),
                SupportedFeatures: optional.EmptyString(),
            },
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
                resp *http.Response
                a AuthorizedNetworkSliceInfo
            )

            a, resp, err := apiClient.NetworkSliceInformationDocumentApi.NSSelectionGet(context.Background(), subtest.nfType, subtest.nfId, subtest.localVarOptionals)

            if err != nil {
                t.Errorf(err.Error())
            }

            if resp.StatusCode != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
            }

            if test.CheckAuthorizedNetworkSliceInfo(a, subtest.expectAuthorizedNetworkSliceInfo) == false {
                e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                r, _ := json.Marshal(a)
                t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }
}
