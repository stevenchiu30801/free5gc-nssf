/*
 * NSSF NS Selection
 *
 * Utility for Client API
 */

package client

import (
    "encoding/json"

    "github.com/antihax/optional"
)

func ConvertQueryParamIntf(intf interface{}) optional.Interface {
    e, _ := json.Marshal(intf)
    return optional.NewInterface(string(e))
}
