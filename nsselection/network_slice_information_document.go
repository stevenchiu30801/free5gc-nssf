/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package nsselection

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "free5gc-nssf/flog"
    . "free5gc-nssf/model"
    "free5gc-nssf/nssf_handler/nssf_message"
    "free5gc-nssf/util"
)

// Parse NSSelectionGet query parameter
func parseQueryParameter(c *gin.Context) (p NsselectionQueryParameter, err error) {

    if c.Query("nf-type") != "" {
        p.NfType = new(NfType)
        *p.NfType = NfType(c.Query("nf-type"))
    }

    p.NfId = c.Query("nf-id")

    if c.Query("slice-info-request-for-registration") != "" {
        p.SliceInfoRequestForRegistration = new(SliceInfoForRegistration)
        err = json.NewDecoder(strings.NewReader(c.Query("slice-info-request-for-registration"))).Decode(p.SliceInfoRequestForRegistration)
        if err != nil {
            return
        }
    }

    if c.Query("slice-info-request-for-pdu-session") != "" {
        p.SliceInfoRequestForPduSession = new(SliceInfoForPduSession)
        err = json.NewDecoder(strings.NewReader(c.Query("slice-info-request-for-pdu-session"))).Decode(p.SliceInfoRequestForPduSession)
        if err != nil {
            return
        }
    }

    if c.Query("home-plmn-id") != "" {
        p.HomePlmnId = new(PlmnId)
        err = json.NewDecoder(strings.NewReader(c.Query("home-plmn-id"))).Decode(p.HomePlmnId)
        if err != nil {
            return
        }
    }

    if c.Query("tai") != "" {
        p.Tai = new(Tai)
        err = json.NewDecoder(strings.NewReader(c.Query("tai"))).Decode(p.Tai)
        if err != nil {
            return
        }
    }

    if c.Query("supported-features") != "" {
        p.SupportedFeatures = c.Query("supported-features")
    }

    return
}

// Check if the NF service consumer is authorized
// TODO: Check if the NF service consumer is legal with local configuration, or possibly after querying NRF through
//       `nf-id` e.g. Whether the V-NSSF is authorized
func checkNfServiceConsumer(nfType NfType) error {
    if nfType != NfType_AMF && nfType != NfType_NSSF {
        return fmt.Errorf("`nf-type`:'%s' is not authorized to retrieve the slice selection information", string(nfType))
    }

    return nil
}

// NSSelectionGet - Retrieve the Network Slice Selection Information
func NSSelectionGet(httpChannel chan nssf_message.HttpResponseMessage, c *gin.Context) {

    flog.Nsselection.Infof("Request received - NSSelectionGet")

    var (
        isValidRequest bool = true
        status int
        a AuthorizedNetworkSliceInfo
        d ProblemDetails
    )

    // TODO: Record request times of the NF service consumer and response with ProblemDetails of 429 Too Many Requests
    //       if the consumer has sent too many requests in a configured amount of time
    // TODO: Check URI length and response with ProblemDetails of 414 URI Too Long if URI is too long

    // Parse query parameter
    p, err := parseQueryParameter(c)
    if err != nil {
        problemDetail := "[Query Parameter] " + err.Error()
        status = http.StatusBadRequest
        d = ProblemDetails {
            Title: util.MALFORMED_REQUEST,
            Status: http.StatusBadRequest,
            Detail: problemDetail,
        }
        isValidRequest = false
    }

    // Check data integrity
    err = p.CheckIntegrity()
    if err != nil {
        problemDetail := "[Query Parameter] " + err.Error()
        s := strings.Split(problemDetail, "`")
        status = http.StatusBadRequest
        if len(s) >= 2 {
            invalidParam := s[len(s) - 2]
            d = ProblemDetails {
                Title: util.INVALID_REQUEST,
                Status: http.StatusBadRequest,
                Detail: problemDetail,
                InvalidParams: []InvalidParam {
                    {
                        Param: invalidParam,
                        Reason: problemDetail,
                    },
                },
            }
        } else {
            d = ProblemDetails {
                Title: util.INVALID_REQUEST,
                Status: http.StatusBadRequest,
                Detail: problemDetail,
            }
        }
        isValidRequest = false
    }

    // Check permission of NF service consumer
    err = checkNfServiceConsumer(*p.NfType)
    if err != nil {
        problemDetail := err.Error()
        status = http.StatusForbidden
        d = ProblemDetails {
            Title: util.UNAUTHORIZED_CONSUMER,
            Status: http.StatusForbidden,
            Detail: problemDetail,
        }
        isValidRequest = false
    }

    if isValidRequest == true {
        if p.SliceInfoRequestForRegistration != nil {
            // Network slice information is requested during the Registration procedure
            status = nsselectionForRegistration(p, &a, &d)
        } else {
            // Network slice information is requested during the PDU session establishment procedure
            status = nsselectionForPduSession(p, &a, &d)
        }
    }


    // Set response
    switch status {
        case http.StatusOK:
            nssf_message.SendHttpResponseMessage(httpChannel, nssf_message.HttpResponseMessageResponse, a)
            flog.Nsselection.Infof("Response code 200 OK")
        case http.StatusBadRequest:
            nssf_message.SendHttpResponseMessage(httpChannel, nssf_message.HttpResponseMessageProblemDetails, d)
            flog.Nsselection.Infof(d.Detail)
            flog.Nsselection.Infof("Response code 400 Bad Request")
        case http.StatusForbidden:
            nssf_message.SendHttpResponseMessage(httpChannel, nssf_message.HttpResponseMessageProblemDetails, d)
            flog.Nsselection.Infof(d.Detail)
            flog.Nsselection.Infof("Response code 403 Forbidden")
        default:
            flog.Nsselection.Warnf("Unknown response code")
    }
}
