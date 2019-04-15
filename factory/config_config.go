/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "errors"
)

type Config struct {

    Info *Info `yaml:"info"`

    Configuration *Configuration `yaml:"configuration"`
}

func (c *Config) checkIntegrity() error {
    if c.Info == nil {
        return errors.New("`info` in configuration should not be empty")
    } else {
        err := c.Info.checkIntegrity()
        if err != nil {
            errMsg := "`info`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if c.Configuration == nil {
        return errors.New("`configuration` in configuration should not be empty")
    } else {
        err := c.Configuration.checkIntegrity()
        if err != nil {
            errMsg := "`configuration`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    return nil
}
