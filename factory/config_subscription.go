/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"
    "time"

    . "../model"
)

type Subscription struct {

    SubscriptionId string `yaml:"subscriptionId"`

    NfNssaiAvailabilityUri string `yaml:"nfNssaiAvailabilityUri"`

    TaiList []Tai `yaml:"taiList"`

    Event *NssfEventType `yaml:"event"`

    Expiry time.Time `yaml:"expiry,omitempty"`
}

func (s *Subscription) checkIntegrity() error {
    if s.SubscriptionId == "" {
        return fmt.Errorf("`subscriptionId` should not be empty")
    }

    if s.NfNssaiAvailabilityUri == "" {
        return fmt.Errorf("`nfNssaiAvailabilityUri` should not be empty")
    }

    if s.TaiList == nil || len(s.TaiList) == 0 {
        return fmt.Errorf("`taiList` should not be empty")
    } else {
        for i, tai := range s.TaiList {
            err := tai.CheckIntegrity()
            if err != nil {
                return fmt.Errorf("`taiList`[%d]:%s", i, err.Error())
            }
        }
    }

    if s.Event == nil {
        return fmt.Errorf("`event` should not be empty")
    } else {
        err := s.Event.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`event`:%s", err.Error())
        }
    }

    return nil
}
