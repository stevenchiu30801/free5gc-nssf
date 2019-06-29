/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package nssf

import (
    "encoding/json"
	"net/http"
    "strings"

    "github.com/gorilla/mux"

    flog "../flog"
    . "../model"
)

// NSSAIAvailabilityDelete - Deletes an already existing S-NSSAIs per TA provided by the NF service consumer (e.g AMF)
func NSSAIAvailabilityDelete(w http.ResponseWriter, r *http.Request) {

    flog.Nssaiavailability.Infof("Request received - NSSAIAvailabilityDelete")

    var (
        nfId string
        status int
        d ProblemDetails
    )

    // Parse nfId from URL path
    vars := mux.Vars(r)
    nfId = vars["nfId"]

    status = nssaiavailabilityDelete(nfId, &d)

    // Set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)
    switch status {
        case http.StatusNoContent:
            flog.Nssaiavailability.Infof("Response code 204 No Content")
        case http.StatusNotFound:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 404 Not Found")
        default:
            flog.Nssaiavailability.Warnf("Unknown response code")
    }
}

// NSSAIAvailabilityPatch - Updates an already existing S-NSSAIs per TA provided by the NF service consumer (e.g AMF)
func NSSAIAvailabilityPatch(w http.ResponseWriter, r *http.Request) {

    flog.Nssaiavailability.Infof("Request received - NSSAIAvailabilityPatch")

    var (
        isValidRequest bool = true
        nfId string
        status int
        p PatchDocument
        a AuthorizedNssaiAvailabilityInfo
        d ProblemDetails
    )

    // Parse nfId from URL path
    vars := mux.Vars(r)
    nfId = vars["nfId"]

    // Parse request body
    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        problemDetail := "[Request Body] " + err.Error()
        status = http.StatusBadRequest
        d = ProblemDetails {
            Title: MALFORMED_REQUEST,
            Status: http.StatusBadRequest,
            Detail: problemDetail,
        }
        isValidRequest = false
    }

    // Check data integrity
    err = p.CheckIntegrity()
    if err != nil {
        problemDetail := "[Request Body] " + err.Error()
        s := strings.Split(problemDetail, "`")
        invalidParam := s[len(s) - 2]
        status = http.StatusBadRequest
        d = ProblemDetails {
            Title: INVALID_REQUEST,
            Status: http.StatusBadRequest,
            Detail: problemDetail,
            InvalidParams: []InvalidParam {
                {
                    Param: invalidParam,
                    Reason: problemDetail,
                },
            },
        }
        isValidRequest = false
    }

    // TODO: Request NfProfile of NfId from NRF
    //       Check if NfId is valid AMF and obtain AMF Set ID
    //       If NfId is invalid, return ProblemDetails with code 404 Not Found
    //       If NF consumer is not authorized to update NSSAI availability, return ProblemDetails with code 403 Forbidden

    if isValidRequest == true {
        status = nssaiavailabilityPatch(nfId, p, &a, &d)
    }

    // Set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)
    switch status {
        case http.StatusOK:
            json.NewEncoder(w).Encode(&a)
            flog.Nssaiavailability.Infof("Response code 200 OK")
        case http.StatusBadRequest:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 400 Bad Request")
        case http.StatusForbidden:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 403 Forbidden")
        case http.StatusNotFound:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 404 Not Found")
        case http.StatusConflict:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 409 Conflict")
        default:
            flog.Nssaiavailability.Warnf("Unknown response code")
    }
}

// NSSAIAvailabilityPut - Updates/replaces the NSSF with the S-NSSAIs the NF service consumer (e.g AMF) supports per TA
func NSSAIAvailabilityPut(w http.ResponseWriter, r *http.Request) {

    flog.Nssaiavailability.Infof("Request received - NSSAIAvailabilityPut")

    var (
        isValidRequest bool = true
        nfId string
        status int
        n NssaiAvailabilityInfo
        a AuthorizedNssaiAvailabilityInfo
        d ProblemDetails
    )

    // Parse nfId from URL path
    vars := mux.Vars(r)
    nfId = vars["nfId"]

    // Parse request body
    err := json.NewDecoder(r.Body).Decode(&n)
    if err != nil {
        problemDetail := "[Request Body] " + err.Error()
        status = http.StatusBadRequest
        d = ProblemDetails {
            Title: MALFORMED_REQUEST,
            Status: http.StatusBadRequest,
            Detail: problemDetail,
        }
        isValidRequest = false
    }

    // Check data integrity
    err = n.CheckIntegrity()
    if err != nil {
        problemDetail := "[Request Body] " + err.Error()
        s := strings.Split(problemDetail, "`")
        invalidParam := s[len(s) - 2]
        status = http.StatusBadRequest
        d = ProblemDetails {
            Title: INVALID_REQUEST,
            Status: http.StatusBadRequest,
            Detail: problemDetail,
            InvalidParams: []InvalidParam {
                {
                    Param: invalidParam,
                    Reason: problemDetail,
                },
            },
        }
        isValidRequest = false
    }

    // TODO: Request NfProfile of NfId from NRF
    //       Check if NfId is valid AMF and obtain AMF Set ID
    //       If NfId is invalid, return ProblemDetails with code 404 Not Found
    //       If NF consumer is not authorized to update NSSAI availability, return ProblemDetails with code 403 Forbidden

    if isValidRequest == true {
        status = nssaiavailabilityPut(nfId, n, &a, &d)
    }

    // Set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)
    switch status {
        case http.StatusOK:
            json.NewEncoder(w).Encode(&a)
            flog.Nssaiavailability.Infof("Response code 200 OK")
        case http.StatusBadRequest:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 400 Bad Request")
        case http.StatusForbidden:
            json.NewEncoder(w).Encode(&d)
            flog.Nssaiavailability.Infof(d.Detail)
            flog.Nssaiavailability.Infof("Response code 403 Forbidden")
        default:
            flog.Nssaiavailability.Warnf("Unknown response code")
    }
}
