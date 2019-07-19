/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssaiavailability

import (
    "context"
    // "encoding/json"
    "net/http"
    "testing"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    // . "free5gc-nssf/model"
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

}

func TestNSSAIAvailabilityPut(t *testing.T) {

}
