/*
 * NSSF Testing Utility
 */

package test

import (
    "reflect"

    . "free5gc-nssf/model"
)

func CheckAuthorizedNetworkSliceInfo(target AuthorizedNetworkSliceInfo, expectList []AuthorizedNetworkSliceInfo) bool {
    target.Sort()
    for _, expectElement := range expectList {
        if reflect.DeepEqual(target, expectElement) == true {
            return true
        }
    }
    return false
}
