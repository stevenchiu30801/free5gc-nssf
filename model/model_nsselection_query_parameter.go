/*
 * NSSF NS Selection
 * 
 * NSSF Network Slice Selection Service
 */

package model

import (
    "fmt"
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
        return fmt.Errorf("`nf-type` should not be empty")
    } else {
        err := p.NfType.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`nf-type`:%s", err.Error())
        }
    }

    if p.NfId == "" {
        return fmt.Errorf("`nf-id` should not be empty")
    }

    if p.SliceInfoRequestForRegistration != nil {
        if p.SliceInfoRequestForPduSession != nil {
            return fmt.Errorf("Slice info requests for both registration and PDU session are provided simultaneously")
        }

        err := p.SliceInfoRequestForRegistration.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`slice-info-request-for-registration`:%s", err.Error())
        }
    } else if p.SliceInfoRequestForPduSession != nil {
        err := p.SliceInfoRequestForPduSession.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`slice-info-request-for-pdu-session`:%s", err.Error())
        }
    } else {
        return fmt.Errorf("None of slice info request for registration or PDU session is provided")
    }

    if p.HomePlmnId != nil {
        err := p.HomePlmnId.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`home-plmn-id`:%s", err.Error())
        }
    }

    if p.Tai != nil {
        err := p.Tai.CheckIntegrity()
        if err != nil {
            return fmt.Errorf("`tai`:%s", err.Error())
        }
    }

    return nil
}
