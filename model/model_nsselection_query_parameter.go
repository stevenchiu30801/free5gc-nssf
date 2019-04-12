/*
 * NSSF NS Selection
 * 
 * NSSF Network Slice Selection Service
 */

package model

import (
    "errors"
)

type NsselectionQueryParameter struct {

    NfType *NfType `json:"nf-type"`

    NfId string `json:"nf-id"`

    SliceInfoRequestForRegistration *SliceInfoForRegistration `json:"slice-info-request-for-registration, omitempty"`

    SliceInfoRequestForPduSession *SliceInfoForPduSession `json:"slice-info-request-for-pdu-session, omitempty"`

    HomePlmnId *PlmnId `json:"home-plmn-id, omitempty"`

    Tai *Tai `json:"tai, omitempty"`

    SupportedFeatures string `json:"supported-features, omitempty"`
}

func (p *NsselectionQueryParameter) CheckIntegrity() error {
    if p.NfType == nil || *p.NfType == "" {
        return errors.New("`nf-type` in query parameter should not be empty")
    } else {
        err := p.NfType.CheckIntegrity()
        if err != nil {
            errMsg := "`nf-type`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if p.NfId == "" {
        return errors.New("`nf-id` in query parameter should not be empty")
    }

    if p.SliceInfoRequestForRegistration != nil {
        if p.SliceInfoRequestForPduSession != nil {
            return errors.New("Slice info requests for both registration and PDU session are provided simultaneously")
        }

        err := p.SliceInfoRequestForRegistration.CheckIntegrity()
        if err != nil {
            errMsg := "`slice-info-request-for-registration`:" + err.Error()
            return errors.New(errMsg)
        }
    } else if p.SliceInfoRequestForPduSession != nil {
        err := p.SliceInfoRequestForPduSession.CheckIntegrity()
        if err != nil {
            errMsg := "`slice-info-request-for-pdu-session`:" + err.Error()
            return errors.New(errMsg)
        }
    } else {
        return errors.New("None of slice info request for registration or PDU session is provided")
    }

    if p.HomePlmnId != nil {
        err := p.HomePlmnId.CheckIntegrity()
        if err != nil {
            errMsg := "`home-plmn-id`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    if p.Tai != nil {
        err := p.Tai.CheckIntegrity()
        if err != nil {
            errMsg := "`tai`:" + err.Error()
            return errors.New(errMsg)
        }
    }

    return nil
}
