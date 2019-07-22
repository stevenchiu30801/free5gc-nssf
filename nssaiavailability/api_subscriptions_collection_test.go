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
    "time"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
    "free5gc-nssf/nssaiavailability/client"
    "free5gc-nssf/test"
)

var testingNssaiavailabilitySubscribeApi = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
}

func TestNSSAIAvailabilityPost(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailabilitySubscribeApi.ConfigFile)
    if testingNssaiavailabilitySubscribeApi.MuteLogInd == true {
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
        nssfEventSubscriptionCreateData *NssfEventSubscriptionCreateData
        expectStatus int
        expectNssfEventSubscriptionCreatedData *NssfEventSubscriptionCreatedData
    }{
        {
            name: "Post",
            nssfEventSubscriptionCreateData: &NssfEventSubscriptionCreateData {
                NfNssaiAvailabilityUri: "http://free5gc-amf2.nctu.me:8080/namf-nssaiavailability/v1/nssai-availability/notify",
                TaiList: []Tai {
                    {
                        PlmnId: &PlmnId {
                            Mcc: "466",
                            Mnc: "92",
                        },
                        Tac: "33456",
                    },
                    {
                        PlmnId: &PlmnId {
                            Mcc: "466",
                            Mnc: "92",
                        },
                        Tac: "33458",
                    },
                },
                Event: func() NssfEventType { n := NssfEventType_SNSSAI_STATUS_CHANGE_REPORT; return n }(),
                Expiry: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2019-06-24T16:35:31+08:00"); return &t }(),
            },
            expectStatus: http.StatusCreated,
            expectNssfEventSubscriptionCreatedData: &NssfEventSubscriptionCreatedData {
                SubscriptionId: "2",
                Expiry: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2019-06-24T16:35:31+08:00"); return &t }(),
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
                            {
                                Sst: 2,
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
            },
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                resp *http.Response
                n NssfEventSubscriptionCreatedData
            )

            n, resp, err := apiClient.SubscriptionsCollectionApi.NSSAIAvailabilityPost(context.Background(), *subtest.nssfEventSubscriptionCreateData)

            if err != nil {
                t.Errorf(err.Error())
            }

            if resp.StatusCode != subtest.expectStatus {
                t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
            }

            if reflect.DeepEqual(n, *subtest.expectNssfEventSubscriptionCreatedData) == false {
                e, _ := json.Marshal(*subtest.expectNssfEventSubscriptionCreatedData)
                r, _ := json.Marshal(n)
                t.Errorf("Incorrect NSSF event subscription created data:\nexpected\n%s\n, got\n%s", string(e), string(r))
            }
        })
    }

    err := srv.Shutdown(context.Background())
    if err != nil {
        t.Fatal(err)
    }
}
