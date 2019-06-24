/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"

    . "../model"
)

type Subscription struct {

    SubscriptionId string `yaml:"subscriptionId"`

    SubscriptionData *NssfEventSubscriptionCreateData `yaml:"subscriptionData"`
}

func (s *Subscription) checkIntegrity() error {
    if s.SubscriptionId == "" {
        return fmt.Errorf("`subscriptionId` should not be empty")
    }

    if s.SubscriptionData == nil {
        return fmt.Errorf("`subscriptionData` should not be empty")
    } else {
        err := s.SubscriptionData.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`subscriptionData`:%s", err.Error())
        }
    }

    return nil
}
