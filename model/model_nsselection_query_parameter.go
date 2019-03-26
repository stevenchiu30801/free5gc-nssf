/*
 * NSSF NS Selection
 * 
 * NSSF Network Slice Selection Service
 */

package model

type NsselectionQueryParameter struct {

    // TODO: use NfType type instead of string type
    // NfType NfType `json:"nf-type"`
    NfType string `json:"nf-type"`

    // NfId NfInstanceId `json:"nf-id"`
    NfId string `json:"nf-id"`

    SliceInfoRequestForRegistration *SliceInfoForRegistration `json:"slice-info-request-for-registration, omitempty"`

    SliceInfoRequestForPduSession *SliceInfoForPduSession `json:"slice-info-request-for-pdu-session, omitempty"`

    HomePlmnId *PlmnId `json:"home-plmn-id, omitempty"`

    Tai *Tai `json:"tai, omitempty"`

    SupportedFeatures string `json:"supported-features, omitempty"`
}
