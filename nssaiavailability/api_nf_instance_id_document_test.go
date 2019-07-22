/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssaiavailability

import (
    "context"
    "encoding/json"
    "net/http"
    "reflect"
    "testing"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
    "free5gc-nssf/nssaiavailability/client"
    "free5gc-nssf/test"
)

var testingNssaiavailabilityStoreApi = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
}

func TestNSSAIAvailabilityDelete(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailabilityStoreApi.ConfigFile)
    if testingNssaiavailabilityStoreApi.MuteLogInd == true {
        flog.Nssaiavailability.MuteLog()
        flog.Util.MuteLog()
    }

    router := NewRouter()
    srv := &http.Server {
        Addr: ":8080",
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            t.Fatal(err)
        }
    }()

    configuration := client.NewConfiguration()
    configuration.SetBasePath("http://localhost:8080")
    apiClient := client.NewAPIClient(configuration)

    subtests := []struct {
        name string
        nfId string
        expectStatus int
    }{
        {
            name: "Delete",
            nfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
            expectStatus: http.StatusNoContent,
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                resp *http.Response
            )

            resp, err := apiClient.NFInstanceIDDocumentApi.NSSAIAvailabilityDelete(context.Background(), subtest.nfId)

            if err != nil {
                t.Errorf(err.Error())
            }

            if resp.StatusCode != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
            }
        })
    }

    err := srv.Shutdown(context.Background())
    if err != nil {
        t.Fatal(err)
    }
}

func TestNSSAIAvailabilityPatch(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailabilityStoreApi.ConfigFile)
    if testingNssaiavailabilityStoreApi.MuteLogInd == true {
        flog.Nssaiavailability.MuteLog()
        flog.Util.MuteLog()
    }

    router := NewRouter()
    srv := &http.Server {
        Addr: ":8080",
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            t.Fatal(err)
        }
    }()

    configuration := client.NewConfiguration()
    configuration.SetBasePath("http://localhost:8080")
    apiClient := client.NewAPIClient(configuration)

    subtests := []struct {
        name string
        nfId string
        patchItem []PatchItem
        expectStatus int
        expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
    }{
        {
            name: "Patch",
            nfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
            patchItem: []PatchItem {
                {
                    Op: func() PatchOperation { p := PatchOperation_ADD; return p }(),
                    Path: "/0/supportedSnssaiList/-",
                    Value: Snssai {
                        Sst: 1,
                        Sd: "1",
                    },
                },
                {
                    Op: func() PatchOperation { p := PatchOperation_COPY; return p }(),
                    Path: "/1/supportedSnssaiList/0",
                    From: "/0/supportedSnssaiList/0",
                },
            },
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
                },
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                resp *http.Response
                a AuthorizedNssaiAvailabilityInfo
            )

            a, resp, err := apiClient.NFInstanceIDDocumentApi.NSSAIAvailabilityPatch(context.Background(), subtest.nfId, subtest.patchItem)

            if err != nil {
                t.Errorf(err.Error())
            }

            if resp.StatusCode != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
            }

            if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
                e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
                r, _ := json.Marshal(a)
                t.Errorf("Incorrect authorized nssai availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }

    err := srv.Shutdown(context.Background())
    if err != nil {
        t.Fatal(err)
    }
}

func TestNSSAIAvailabilityPut(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailabilityStoreApi.ConfigFile)
    if testingNssaiavailabilityStoreApi.MuteLogInd == true {
        flog.Nssaiavailability.MuteLog()
        flog.Util.MuteLog()
    }

    router := NewRouter()
    srv := &http.Server {
        Addr: ":8080",
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            t.Fatal(err)
        }
    }()

    configuration := client.NewConfiguration()
    configuration.SetBasePath("http://localhost:8080")
    apiClient := client.NewAPIClient(configuration)

    subtests := []struct {
        name string
        nfId string
        nssaiAvailabilityInfo NssaiAvailabilityInfo
        expectStatus int
        expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
    }{
        {
            name: "Put",
            nfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
            nssaiAvailabilityInfo: NssaiAvailabilityInfo {
                SupportedNssaiAvailabilityData: []SupportedNssaiAvailabilityData {
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
                    },
                },
                SupportedFeatures: "",
            },
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
                resp *http.Response
                a AuthorizedNssaiAvailabilityInfo
            )

            a, resp, err := apiClient.NFInstanceIDDocumentApi.NSSAIAvailabilityPut(context.Background(), subtest.nfId, subtest.nssaiAvailabilityInfo)

            if err != nil {
                t.Errorf(err.Error())
            }

            if resp.StatusCode != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
            }

            if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
                e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
                r, _ := json.Marshal(a)
                t.Errorf("Incorrect authorized nssai availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }

    err := srv.Shutdown(context.Background())
    if err != nil {
        t.Fatal(err)
    }
}
