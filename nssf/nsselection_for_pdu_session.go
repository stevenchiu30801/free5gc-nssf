/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "net/http"

    . "../model"
)

func nsselectionForPduSession(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo, d *ProblemDetails) int {

    return http.StatusOK
}
