/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "net/http"

    flog "../flog"
    . "../model"
)

// Set Allowed NSSAI with Subscribed S-NSSAI(s) which are marked as default S-NSSAI(s)
func useDefaultSubscribedSnssai(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo) {
    var mappingOfSnssai []MappingOfSnssai
    if p.HomePlmnId != nil {
        // Find mapping of Subscribed S-NSSAI of UE's HPLMN to S-NSSAI in Serving PLMN from NSSF configuration
        mappingOfSnssai = getMappingOfPlmnFromConfig(*p.HomePlmnId)

        if mappingOfSnssai == nil {
            flog.Warnf("No S-NSSAI mapping of UE's HPLMN %+v in NSSF configuration", *p.HomePlmnId)
            return
        }
    }

    for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
        if subscribedSnssai.DefaultIndication == true {
            // Subscribed S-NSSAI is marked as default S-NSSAI

            var mappingOfSubscribedSnssai Snssai
            if p.HomePlmnId != nil && checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == false {
                targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai, mappingOfSnssai)

                if found == false {
                    flog.Warnf("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
                              *subscribedSnssai.SubscribedSnssai,
                              *p.HomePlmnId)
                    continue
                } else {
                    mappingOfSubscribedSnssai = *targetMapping.ServingSnssai
                }
            } else {
                mappingOfSubscribedSnssai = *subscribedSnssai.SubscribedSnssai
            }

            if p.Tai != nil && checkSupportedSnssaiInTa(mappingOfSubscribedSnssai, p.Tai.Tac) == false {
                continue
            }

            var allowedSnssaiElement AllowedSnssai
            allowedSnssaiElement.AllowedSnssai = new(Snssai)
            *allowedSnssaiElement.AllowedSnssai = mappingOfSubscribedSnssai
            nsiInformationList := getNsiInformationListFromConfig(mappingOfSubscribedSnssai)
            if nsiInformationList != nil {
                allowedSnssaiElement.NsiInformationList = append(allowedSnssaiElement.NsiInformationList,
                                                                 nsiInformationList...)
            }
            if p.HomePlmnId != nil && checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == false {
                allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
            }

            // Default Access Type is set to 3GPP Access if no TAI is provided
            // TODO: Depend on operator implementation, it may also return S-NSSAIs in all valid Access Type if
            //       UE's Access Type could not be identified
            var accessType AccessType = IS_3_GPP_ACCESS
            if p.Tai != nil {
                accessType = getAccessTypeFromConfig(*p.Tai)
            }

            addAllowedSnssai(allowedSnssaiElement, accessType, a)

        }
    }

    addAmfInformation(a)
}

// Set Configured NSSAI with S-NSSAI(s) in Requested NSSAI which are marked as Default Configured NSSAI
func useDefaultConfiguredNssai(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo) {
    for _, requestedSnssai := range p.SliceInfoRequestForRegistration.RequestedNssai {
        // Check whether the Default Configured S-NSSAI is standard, which could be commonly decided by all roaming partners
        if checkStandardSnssai(requestedSnssai) == false {
            flog.Infof("S-NSSAI %+v in Requested NSSAI which based on Default Configured NSSAI is not standard", requestedSnssai)
            continue
        }

        // Check whether the Default Configured S-NSSAI is subscribed
        for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
            if requestedSnssai == *subscribedSnssai.SubscribedSnssai {
                var configuredSnssai ConfiguredSnssai
                configuredSnssai.ConfiguredSnssai = new(Snssai)
                *configuredSnssai.ConfiguredSnssai = requestedSnssai

                a.ConfiguredNssai = append(a.ConfiguredNssai, configuredSnssai)
                break
            }
        }
    }
}

// Set Configured NSSAI with Subscribed S-NSSAI(s)
func setConfiguredNssai(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo) {
    var mappingOfSnssai []MappingOfSnssai
    if p.HomePlmnId != nil {
        // Find mapping of Subscribed S-NSSAI of UE's HPLMN to S-NSSAI in Serving PLMN from NSSF configuration
        mappingOfSnssai = getMappingOfPlmnFromConfig(*p.HomePlmnId)

        if mappingOfSnssai == nil {
            flog.Warnf("No S-NSSAI mapping of UE's HPLMN %+v in NSSF configuration", *p.HomePlmnId)
            return
        }
    }

    for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
        var mappingOfSubscribedSnssai Snssai
        if p.HomePlmnId != nil && checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == false {
            targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai, mappingOfSnssai)

            if found == false {
                flog.Warnf("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
                          *subscribedSnssai.SubscribedSnssai,
                          *p.HomePlmnId)
                continue
            } else {
                mappingOfSubscribedSnssai = *targetMapping.ServingSnssai
            }
        } else {
            mappingOfSubscribedSnssai = *subscribedSnssai.SubscribedSnssai
        }

        if checkSupportedSnssaiInPlmn(mappingOfSubscribedSnssai) == true {
            var configuredSnssai ConfiguredSnssai
            configuredSnssai.ConfiguredSnssai = new(Snssai)
            *configuredSnssai.ConfiguredSnssai = mappingOfSubscribedSnssai
            if p.HomePlmnId != nil && checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == false {
                configuredSnssai.MappedHomeSnssai = new(Snssai)
                *configuredSnssai.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
            }

            a.ConfiguredNssai = append(a.ConfiguredNssai, configuredSnssai)
        }
    }
}

// Network slice selection for registration
// The function is executed when the IE, `slice-info-request-for-registration`, is provided in query parameters
func nsselectionForRegistration(p NsselectionQueryParameter,
                                a *AuthorizedNetworkSliceInfo,
                                d *ProblemDetails) (status int) {
    if p.HomePlmnId != nil {
        // Check whether UE's Home PLMN is supported when UE is a roamer
        if checkSupportedHplmn(*p.HomePlmnId) == false {
            for _, requestedSnssai := range p.SliceInfoRequestForRegistration.RequestedNssai {
                a.RejectedNssaiInPlmn = append(a.RejectedNssaiInPlmn, requestedSnssai)
            }

            status = http.StatusOK
            return
        }
    }

    if p.Tai != nil {
        // Check whether UE's current TA is supported when UE provides TAI
        if checkSupportedTa(*p.Tai) == false {
            for _, requestedSnssai := range p.SliceInfoRequestForRegistration.RequestedNssai {
                a.RejectedNssaiInTa = append(a.RejectedNssaiInTa, requestedSnssai)
            }

            status = http.StatusOK
            return
        }
    }

    if p.SliceInfoRequestForRegistration.RequestMapping == true {
        // Based on TS 29.531 v15.2.0, when `requestMapping` is set to true, the NSSF shall return the VPLMN specific
        // mapped S-NSSAI values for the S-NSSAI values in `subscribedNssai`. But also `sNssaiForMapping` shall be
        // provided if `requestMapping` is set to true. In the implementation, the NSSF would return mapped S-NSSAIs
        // for S-NSSAIs in both `sNssaiForMapping` and `subscribedSnssai` if present

        if p.HomePlmnId == nil {
            problemDetail := "`home-plmn-id` should be provided when requesting VPLMN specific mapped S-NSSAI values"
            *d = ProblemDetails {
                Title: INVALID_REQUEST,
                Status: http.StatusBadRequest,
                Detail: problemDetail,
                InvalidParams: []InvalidParam {
                    {
                        Param: "home-plmn-id",
                        Reason: problemDetail,
                    },
                },
            }

            status = http.StatusBadRequest
            return
        }

        mappingOfSnssai := getMappingOfPlmnFromConfig(*p.HomePlmnId)

        if mappingOfSnssai != nil {
            // Find mappings for S-NSSAIs in `subscribedSnssai`
            for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
                if checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == true {
                    continue
                }

                targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai, mappingOfSnssai)

                if found == false {
                    flog.Warnf("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
                              *subscribedSnssai.SubscribedSnssai,
                              *p.HomePlmnId)
                    continue
                } else {
                    // Add mappings to Allowed NSSAI list
                    var allowedSnssaiElement AllowedSnssai
                    allowedSnssaiElement.AllowedSnssai = new(Snssai)
                    *allowedSnssaiElement.AllowedSnssai = *targetMapping.ServingSnssai
                    allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                    *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai

                    // Default Access Type is set to 3GPP Access if no TAI is provided
                    // TODO: Depend on operator implementation, it may also return S-NSSAIs in all valid Access Type if
                    //       UE's Access Type could not be identified
                    var accessType AccessType = IS_3_GPP_ACCESS
                    if p.Tai != nil {
                        accessType = getAccessTypeFromConfig(*p.Tai)
                    }

                    addAllowedSnssai(allowedSnssaiElement, accessType, a)
                }
            }

            // Find mappings for S-NSSAIs in `sNssaiForMapping`
            for _, snssai := range p.SliceInfoRequestForRegistration.SNssaiForMapping {
                if checkStandardSnssai(snssai) == true {
                    continue
                }

                targetMapping, found := findMappingWithHomeSnssai(snssai, mappingOfSnssai)

                if found == false {
                    flog.Warnf("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
                              snssai,
                              *p.HomePlmnId)
                    continue
                } else {
                    // Add mappings to Allowed NSSAI list
                    var allowedSnssaiElement AllowedSnssai
                    allowedSnssaiElement.AllowedSnssai = new(Snssai)
                    *allowedSnssaiElement.AllowedSnssai = *targetMapping.ServingSnssai
                    allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                    *allowedSnssaiElement.MappedHomeSnssai = snssai

                    // Default Access Type is set to 3GPP Access if no TAI is provided
                    // TODO: Depend on operator implementation, it may also return S-NSSAIs in all valid Access Type if
                    //       UE's Access Type could not be identified
                    var accessType AccessType = IS_3_GPP_ACCESS
                    if p.Tai != nil {
                        accessType = getAccessTypeFromConfig(*p.Tai)
                    }

                    addAllowedSnssai(allowedSnssaiElement, accessType, a)
                }
            }

            status = http.StatusOK
            return
        } else {
            flog.Warnf("No S-NSSAI mapping of UE's HPLMN %+v in NSSF configuration", *p.HomePlmnId)

            status = http.StatusOK
            return
        }
    }

    checkInvalidRequestedNssai := false
    if p.SliceInfoRequestForRegistration.RequestedNssai != nil && len(p.SliceInfoRequestForRegistration.RequestedNssai) != 0 {
        // Requested NSSAI is provided
        // Verify which S-NSSAI(s) in the Requested NSSAI are permitted based on comparing the Subscribed S-NSSAI(s)

        if checkSupportedNssaiInPlmn(p.SliceInfoRequestForRegistration.RequestedNssai) == false {
            // Return ProblemDetails indicating S-NSSAI is not supported
            // TODO: Based on TS 23.501 V15.2.0, if the Requested NSSAI includes an S-NSSAI that is not valid in the
            //       Serving PLMN, the NSSF may derive the Configured NSSAI for Serving PLMN
            *d = ProblemDetails {
                Title: UNSUPPORTED_RESOURCE,
                Status: http.StatusForbidden,
                Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
                Cause: "SNSSAI_NOT_SUPPORTED",
            }

            status = http.StatusForbidden
            return
        }

        // Check if any Requested S-NSSAIs is present in Subscribed S-NSSAIs
        checkIfRequestAllowed := false

        for _, requestedSnssai := range p.SliceInfoRequestForRegistration.RequestedNssai {
            if p.Tai != nil && checkSupportedSnssaiInTa(requestedSnssai, p.Tai.Tac) == false {
                // Requested S-NSSAI does not supported in UE's current TA
                // Add it to Rejected NSSAI in TA
                a.RejectedNssaiInTa = append(a.RejectedNssaiInTa, requestedSnssai)
                continue
            }

            var mappingOfRequestedSnssai Snssai
            if p.HomePlmnId != nil && checkStandardSnssai(requestedSnssai) == false {
                // Standard S-NSSAIs are supported to be commonly decided by all roaming partners
                // Only non-standard S-NSSAIs are required to find mappings
                targetMapping, found := findMappingWithServingSnssai(requestedSnssai,
                                                                     p.SliceInfoRequestForRegistration.MappingOfNssai)

                if found == false {
                    // No mapping of Requested S-NSSAI to HPLMN S-NSSAI is provided by UE
                    // TODO: Search for local configuration if there is no provided mapping from UE
                    checkInvalidRequestedNssai = true
                    a.RejectedNssaiInPlmn = append(a.RejectedNssaiInPlmn, requestedSnssai)
                    continue
                } else {
                    mappingOfRequestedSnssai = *targetMapping.HomeSnssai
                }
            } else {
                mappingOfRequestedSnssai = requestedSnssai
            }

            hitSubscription := false
            for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
                if mappingOfRequestedSnssai == *subscribedSnssai.SubscribedSnssai {
                    // Requested S-NSSAI matches one of Subscribed S-NSSAI
                    // Add it to Allowed NSSAI list
                    hitSubscription = true

                    var allowedSnssaiElement AllowedSnssai
                    allowedSnssaiElement.AllowedSnssai = new(Snssai)
                    *allowedSnssaiElement.AllowedSnssai = requestedSnssai
                    nsiInformationList := getNsiInformationListFromConfig(requestedSnssai)
                    if nsiInformationList != nil {
                        allowedSnssaiElement.NsiInformationList = append(allowedSnssaiElement.NsiInformationList,
                                                                         nsiInformationList...)
                    }
                    if p.HomePlmnId != nil && checkStandardSnssai(requestedSnssai) == false {
                        allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                        *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
                    }

                    // Default Access Type is set to 3GPP Access if no TAI is provided
                    // TODO: Depend on operator implementation, it may also return S-NSSAIs in all valid Access Type if
                    //       UE's Access Type could not be identified
                    var accessType AccessType = IS_3_GPP_ACCESS
                    if p.Tai != nil {
                        accessType = getAccessTypeFromConfig(*p.Tai)
                    }

                    addAllowedSnssai(allowedSnssaiElement, accessType, a)

                    checkIfRequestAllowed = true
                    break
                }
            }

            if hitSubscription == false {
                // Requested S-NSSAI does not match any Subscribed S-NSSAI
                // Add it to Rejected NSSAI in PLMN
                checkInvalidRequestedNssai = true
                a.RejectedNssaiInPlmn = append(a.RejectedNssaiInPlmn, requestedSnssai)
            }
        }

        if checkIfRequestAllowed == true {
            addAmfInformation(a)
        } else {
            // No S-NSSAI from Requested NSSAI is present in Subscribed S-NSSAIs
            // Subscribed S-NSSAIs marked as default are used
            useDefaultSubscribedSnssai(p, a)
        }
    } else {
        // No Requested NSSAI is provided
        // Subscribed S-NSSAIs marked as default are used
        checkInvalidRequestedNssai = true
        useDefaultSubscribedSnssai(p, a)
    }

    if p.SliceInfoRequestForRegistration.DefaultConfiguredSnssaiInd == true {
        // Default Configured NSSAI Indication is received from AMF
        // Determine the Configured NSSAI based on the Default Configured NSSAI
        useDefaultConfiguredNssai(p, a)
    } else if checkInvalidRequestedNssai == true {
        // No Requested NSSAI is provided or the Requested NSSAI includes an S-NSSAI that is not valid
        // Determine the Configured NSSAI based on the subscription
        setConfiguredNssai(p, a)
    }

    status = http.StatusOK
    return
}
