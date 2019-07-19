/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssaiavailability

import (
    "bytes"
    "context"
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

var testingNotification = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
    NfNssaiAvailabilityUri: "http://localhost:8080/notification",
}

func generateNotificationRequest() NssfEventNotification {
    const jsonRequest = `
        {
            "subscriptionId": "1",
            "authorizedNssaiAvailabilityData": [
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
                        },
                        {
                            "sst": 2
                        }
                    ]
                }
            ]
        }
    `

    var n NssfEventNotification
    json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n)

    return n
}

func TestNotificationTemplate(t *testing.T) {
    t.Skip()

    // Tests may have different configuration files
    factory.InitConfigFactory(testingNotification.ConfigFile)

    d, _ := yaml.Marshal(*factory.NssfConfig.Info)
    t.Logf("%s", string(d))
}

func TestNotificationPost(t *testing.T) {
    factory.InitConfigFactory(testingNotification.ConfigFile)
    if testingNotification.MuteLogInd == true {
        flog.Nssaiavailability.MuteLog()
    }

    subtests := []struct {
        name string
        generateRequestBody func() NssfEventNotification
        expectRequestBody string
        expectError error
    }{
        {
            name: "Notify",
            generateRequestBody: generateNotificationRequest,
            expectRequestBody: `{"subscriptionId":"1","authorizedNssaiAvailabilityData":[{` +
                `"tai":{"plmnId":{"mcc":"466","mnc":"92"},"tac":"33456"},` +
                `"supportedSnssaiList":[{"sst":1},{"sst":1,"sd":"1"},{"sst":1,"sd":"2"},{"sst":2}]}]}`,
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                n NssfEventNotification
                requestBody string
                returnedError error
            )

            // Create a server to accept testing requests
            mux := http.NewServeMux()
            server := http.Server{Addr: ":8080", Handler: mux}
            ctx, cancel := context.WithCancel(context.Background())
            mux.HandleFunc("/notification", func(w http.ResponseWriter, r *http.Request) {
                buf := new(bytes.Buffer)
                buf.ReadFrom(r.Body)
                requestBody = buf.String()

                w.Header().Set("Content-Type", "application/json; charset=UTF-8")
                w.WriteHeader(http.StatusNoContent)
                cancel()
            })

            // Run the server in a goroutine
            go func() {
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                    t.Fatal(err)
                }
            }()

            // Start to generate and send notification request after channel is closed
            if subtest.generateRequestBody != nil {
                n = subtest.generateRequestBody()
            }

            returnedError = notificationPost(testingNotification.NfNssaiAvailabilityUri, n)

            select {
                case <-ctx.Done():
                    // Shutdown server after accepting the first request
                    server.Shutdown(ctx)
            }

            if returnedError == nil {
                if requestBody != subtest.expectRequestBody {
                    t.Errorf("Incorrect request body:\nexpected\n%s\n, got\n%s", subtest.expectRequestBody, requestBody)
                }
            } else {
                if reflect.DeepEqual(returnedError, subtest.expectError) == false {
                    t.Errorf("Incorrect returned error:\nexpected\n%s\n, got\n%s", subtest.expectError.Error(), returnedError.Error())
                }
            }
        })
    }
}
