/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "encoding/json"
    "reflect"
    "testing"

    "gopkg.in/yaml.v2"

    factory "../factory"
    flog "../flog"
    . "../model"
)

const (
    utilsConfig string = "../test/conf/test_nssf_config.yaml"
    utilsMuteLog bool = false
)

func TestPluginTemplate(t *testing.T) {
    t.Skip()

    factory.InitConfigFactory(utilsConfig)

    d, _ := yaml.Marshal(&factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestAddAmfInformation(t *testing.T) {
    factory.InitConfigFactory(utilsConfig)
    if utilsMuteLog == true {
        flog.Nsselection.MuteLog()
    }

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
                                        NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8085/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8085/nnrf-nfm/v1/nf-instances",
                                        NsiId: "3",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                CandidateAmfList: []string {
                    "0e8831c3-6286-4689-ab27-1e2161e15cb1",
                    "a1fba9ba-2e39-4e22-9c74-f749da571d0d",
                    "ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
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
                                        NrfId: "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
                                    },
                                },
                            },
                        },
                        AccessType: func() *AccessType { a := IS_3_GPP_ACCESS; return &a }(),
                    },
                },
                TargetAmfSet: "2",
                NrfAmfSet: "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            a := *subtest.authorizedNetworkSliceInfo

            addAmfInformation(&a)

            a.Sort()
            if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) != true {
                e, _ := json.Marshal(subtest.expectAuthorizedNetworkSliceInfo)
                r, _ := json.Marshal(&a)
                t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }
}
