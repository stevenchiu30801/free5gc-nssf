/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf

import (
    "encoding/json"
    "net/http"
    "reflect"
    "strings"
    "testing"
    "time"

    "gopkg.in/yaml.v2"

    factory "../factory"
    flog "../flog"
    . "../model"
    test "../test"
)

var testingSubscription = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
}

func generatePostRequest() NssfEventSubscriptionCreateData {
    const jsonRequest = `
        {
            "nfNssaiAvailabilityUri": "free5gc-amf.nctu.me:8080/nnssf-nssaiavailability/v1/nssai-availability/notify",
            "taiList": [
                {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                },
                {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33457"
                }
            ],
            "event": "SNSSAI_STATUS_CHANGE_REPORT",
            "expiry": "2019-06-24T16:35:31+08:00"
        }
    `

    var n NssfEventSubscriptionCreateData
    json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n)

    return n
}

func TestSubscriptionTemplate(t *testing.T) {
    t.Skip()

    // Tests may have different configuration files
    factory.InitConfigFactory(testingSubscription.ConfigFile)

    d, _ := yaml.Marshal(*factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestSubscriptionPost(t *testing.T) {
    factory.InitConfigFactory(testingSubscription.ConfigFile)
    if testingSubscription.MuteLogInd == true {
        flog.Nsselection.MuteLog()
    }

    subtests := []struct {
        name string
        generateRequestBody func() NssfEventSubscriptionCreateData
        expectStatus int
        expectSubscriptionCreatedData *NssfEventSubscriptionCreatedData
        expectProblemDetails *ProblemDetails
    }{
        {
            name: "Subscribe",
            generateRequestBody: generatePostRequest,
            expectStatus: http.StatusCreated,
            expectSubscriptionCreatedData: &NssfEventSubscriptionCreatedData {
                SubscriptionId: "1",
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
                n NssfEventSubscriptionCreateData
                status int
                c NssfEventSubscriptionCreatedData
                d ProblemDetails
            )

            if subtest.generateRequestBody != nil {
                n = subtest.generateRequestBody()
            }

            status = subscriptionPost(n, &c, &d)

            if status == http.StatusCreated {
                if reflect.DeepEqual(c, *subtest.expectSubscriptionCreatedData) == false {
                    e, _ := json.Marshal(*subtest.expectSubscriptionCreatedData)
                    r, _ := json.Marshal(c)
                    t.Errorf("Incorrect NSSF event subscription created data:\nexpected\n%s\n, got\n%s", string(e), string(r))
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
