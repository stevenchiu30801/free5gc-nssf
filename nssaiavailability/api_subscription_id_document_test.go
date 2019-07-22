/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssaiavailability

import (
    "context"
    "net/http"
    "testing"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    "free5gc-nssf/nssaiavailability/client"
    "free5gc-nssf/test"
)

var testingNssaiavailabilityUnsubscribeApi = test.TestingNssaiavailability {
    ConfigFile: test.ConfigFileFromArgs,
    MuteLogInd: test.MuteLogIndFromArgs,
}

func TestNSSAIAvailabilityUnsubscribe(t *testing.T) {
    factory.InitConfigFactory(testingNssaiavailabilityUnsubscribeApi.ConfigFile)
    if testingNssaiavailabilityUnsubscribeApi.MuteLogInd == true {
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
        subscriptionId string
        expectStatus int
    }{
        {
            name: "Delete",
            subscriptionId: "3",
            expectStatus: http.StatusNoContent,
        },
    }

    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            var (
                resp *http.Response
            )

            resp, err := apiClient.SubscriptionIDDocumentApi.NSSAIAvailabilityUnsubscribe(context.Background(), subtest.subscriptionId)

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
