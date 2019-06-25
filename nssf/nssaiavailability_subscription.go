/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf

import (
    "fmt"
    "math"
    "net/http"
    "strconv"
    "time"

    factory "../factory"
    flog "../flog"
    . "../model"
)

// Get available subscription ID from configuration
// In this implementation, string converted from 32-bit integer is used as subscription ID
func getUnusedSubscriptionId() (string, error) {
    var idx uint32 = 1
    for _, subscription := range factory.NssfConfig.Subscriptions {
        tempId, _ := strconv.Atoi(subscription.SubscriptionId)
        if uint32(tempId) == idx {
            if idx == math.MaxUint32 {
                return "", fmt.Errorf("No available subscription ID")
            }
            idx = idx + 1
        } else {
            break
        }
    }
    return strconv.Itoa(int(idx)), nil
}

// NSSAIAvailability subscription POST method
func subscriptionPost(n NssfEventSubscriptionCreateData, c *NssfEventSubscriptionCreatedData, d *ProblemDetails) (status int) {
    var subscription factory.Subscription
    tempId, err := getUnusedSubscriptionId()
    if err != nil {
        flog.Nssaiavailability.Warnf(err.Error())

        *d = ProblemDetails {
            Title: UNSUPPORTED_RESOURCE,
            Status: http.StatusNotFound,
            Detail: err.Error(),
        }

        status = http.StatusNotFound
        return
    }

    subscription.SubscriptionId = tempId
    subscription.SubscriptionData = new(NssfEventSubscriptionCreateData)
    *subscription.SubscriptionData = n

    factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions, subscription)

    c.SubscriptionId = subscription.SubscriptionId
    if subscription.SubscriptionData.Expiry.IsZero() == false {
        c.Expiry = new(time.Time)
        *c.Expiry = subscription.SubscriptionData.Expiry
    }
    c.AuthorizedNssaiAvailabilityData = authorizeOfTaListFromConfig(subscription.SubscriptionData.TaiList)

    status = http.StatusCreated
    return
}

func subscriptionDelete(subscriptionId string, d *ProblemDetails) (status int) {
    for i, subscription := range factory.NssfConfig.Subscriptions {
        if subscription.SubscriptionId == subscriptionId {
            factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions[:i],
                                                      factory.NssfConfig.Subscriptions[i + 1:]...)

            status = http.StatusNoContent
            return
        }
    }

    // No specific subscription ID exists
    problemDetail := fmt.Sprintf("Subscription ID '%s' is not available", subscriptionId)
    *d = ProblemDetails {
        Title: UNSUPPORTED_RESOURCE,
        Status: http.StatusNotFound,
        Detail: problemDetail,
    }

    status = http.StatusNotFound
    return
}
