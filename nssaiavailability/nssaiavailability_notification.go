/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssaiavailability

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"

    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
)

func notificationPost(url string, n NssfEventNotification) error {

    e, _ := json.Marshal(n)

    r, err := http.Post(url, "application/json; charset=UTF-8", bytes.NewReader(e))
    if err != nil {
        return err
    }

    flog.Nssaiavailability.Infof("Request sent - NotificationPost")

    if r.StatusCode == http.StatusNoContent {
        return nil
    } else {
        return fmt.Errorf("Incorrect HTTP status code %s", r.Status)
    }
}
