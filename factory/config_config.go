/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"
)

type Config struct {

    Info *Info `yaml:"info"`

    Configuration *Configuration `yaml:"configuration"`

    Subscriptions []Subscription `yaml:"subscriptions,omitempty"`
}

func (c *Config) CheckIntegrity() error {
    if c.Info == nil {
        return fmt.Errorf("`info` should not be empty")
    } else {
        err := c.Info.checkIntegrity()
        if err != nil {
            return fmt.Errorf("`info`:%s", err.Error())
        }
    }

    if c.Configuration == nil {
        return fmt.Errorf("`configuration` should not be empty")
    } else {
        err := c.Configuration.checkIntegrity()
        if err != nil {
            return fmt.Errorf("`configuration`:%s", err.Error())
        }
    }

    if c.Subscriptions != nil && len(c.Subscriptions) != 0 {
        for i, subscription := range c.Subscriptions {
            err := subscription.checkIntegrity()
            if err != nil {
                return fmt.Errorf("`subscriptions`[%d]:%s", i, err.Error())
            }
        }
    }

    return nil
}
