/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "net/http"

    factory "../factory"
    flog "../flog"
    . "../model"
)

// Check whether UE's Home PLMN is configured/supported
func checkSupportedHplmn(homePlmnId PlmnId) bool {
    for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
        if *mappingFromPlmn.HomePlmnId == homePlmnId {
            return true
        }
    }
    flog.Warn("No Home PLMN %+v in NSSF configuration", homePlmnId)
    return false
}

// Check whether UE's current TA is configured/supported
func checkSupportedTa(tai Tai) bool {
    if *tai.PlmnId != *factory.NssfConfig.Info.ServingPlmnId {
        flog.Info("Invalid PLMN ID %+v provided in TAI", *tai.PlmnId)
        return false
    }

    for _, taConfig := range factory.NssfConfig.Configuration.TaList {
        if taConfig.Tac == tai.Tac {
            return true
        }
    }
    flog.Warn("No TA {Tac:%s} in NSSF configuration", tai.Tac)
    return false
}

// Check whether S-NSSAIs in NSSAI are supported or not in PLMN
func checkSupportedNssaiInPlmn(nssai []Snssai) bool {
    for _, snssai := range nssai {
        // Standard S-NSSAIs are supposed to be supported
        // If not, disable following check and be sure to add supported standard S-NSSAI(s) in configuration
        if checkStandardSnssai(snssai) == true {
            continue
        }

        hitSupportedNssai := false
        for _, supportedSnssai := range factory.NssfConfig.Configuration.SupportedNssaiInPlmn {
            if snssai == supportedSnssai {
                hitSupportedNssai = true
                break
            }
        }

        if hitSupportedNssai == false {
            return false
        }
    }
    return true
}

// Check whether S-NSSAI is supported or not in UE's current TA
func checkSupportedSnssaiInTa(snssai Snssai, tac string) bool {
    for _, taConfig := range factory.NssfConfig.Configuration.TaList {
        if taConfig.Tac == tac {
            for _, supportedSnssai := range taConfig.SupportedNssai {
                if snssai == supportedSnssai {
                    return true
                }
            }
            return false
        }
    }
    flog.Warn("No TA %s in NSSF configuration", tac)
    return false
}

// Check whether S-NSSAI is standard or non-standard value
// A standard S-NSSAI is only comprised of a standardized SST value and no SD
func checkStandardSnssai(snssai Snssai) bool {
    if snssai.Sst >= 1 && snssai.Sst <= 3 && snssai.Sd == "" {
        return true
    }
    return false
}

// Check whether the NSSAI contains the specific S-NSSAI
func checkSnssaiInNssai(targetSnssai Snssai, nssai []Snssai) bool {
    for _, snssai := range nssai {
        if snssai == targetSnssai {
            return true
        }
    }
    return false
}

// Get S-NSSAI mappings of the given Home PLMN ID from configuration
func getMappingOfPlmnFromConfig(homePlmnId PlmnId) []MappingOfSnssai {
    for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
        if *mappingFromPlmn.HomePlmnId == homePlmnId {
            return mappingFromPlmn.MappingOfSnssai
        }
    }
    return nil
}

// Get NSI information list of the given S-NSSAI from configuration
func getNsiInformationListFromConfig(snssai Snssai) []NsiInformation {
    for _, nsiConfig := range factory.NssfConfig.Configuration.NsiList {
        if *nsiConfig.Snssai == snssai {
            return nsiConfig.NsiInformationList
        }
    }
    return nil
}

// Get Access Type of the given TAI from configuraion
func getAccessTypeFromConfig(tai Tai) AccessType {
    for _, taConfig := range factory.NssfConfig.Configuration.TaList {
        if taConfig.Tac == tai.Tac {
            return *taConfig.AccessType
        }
    }
    flog.Warn("No TA {Tac:%s} in NSSF configuration", tai.Tac)
    return IS_3_GPP_ACCESS
}

// Find target S-NSSAI mapping with serving S-NSSAIs from mapping of S-NSSAI(s)
func findMappingWithServingSnssai(snssai Snssai, mappings []MappingOfSnssai) (MappingOfSnssai, bool) {
    for _, mapping := range mappings {
        if *mapping.ServingSnssai == snssai {
            return mapping, true
        }
    }
    return MappingOfSnssai{}, false
}

// Find target S-NSSAI mapping with home S-NSSAIs from mapping of S-NSSAI(s)
func findMappingWithHomeSnssai(snssai Snssai, mappings []MappingOfSnssai) (MappingOfSnssai, bool) {
    for _, mapping := range mappings {
        if *mapping.HomeSnssai == snssai {
            return mapping, true
        }
    }
    return MappingOfSnssai{}, false
}

// Add Allowed S-NSSAI to Authorized Network Slice Info
func addAllowedSnssai(allowedSnssai AllowedSnssai, accessType AccessType, a *AuthorizedNetworkSliceInfo) {
    hitAllowedNssai := false
    for i := range a.AllowedNssaiList {
        if *a.AllowedNssaiList[i].AccessType == accessType {
            hitAllowedNssai = true
            if len(a.AllowedNssaiList[i].AllowedSnssaiList) == 8 {
                flog.Info("Unable to add a new Allowed S-NSSAI since already eight S-NSSAIs in Allowed NSSAI")
            } else {
                a.AllowedNssaiList[i].AllowedSnssaiList = append(a.AllowedNssaiList[i].AllowedSnssaiList, allowedSnssai)
            }
            break
        }
    }

    if hitAllowedNssai == false {
        var allowedNssaiElement AllowedNssai
        allowedNssaiElement.AllowedSnssaiList = append(allowedNssaiElement.AllowedSnssaiList, allowedSnssai)
        allowedNssaiElement.AccessType = &accessType

        a.AllowedNssaiList = append(a.AllowedNssaiList, allowedNssaiElement)
    }
}

// Add AMF information to Authorized Network Slice Info
func addAmfInformation(a *AuthorizedNetworkSliceInfo) {
    if a.AllowedNssaiList == nil || len(a.AllowedNssaiList) == 0 {
        return
    }

    // Check if any AMF can serve the UE
    // That is, whether NSSAI of all Allowed S-NSSAIs is a subset of NSSAI supported by AMF
    // Simply use the first applicable AMF set
    // TODO: Policies of AMF selection (e.g. load balance between AMF instances)
    hitAmfSet := false
    for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
        hitAllowedNssai := true
        for _, allowedNssai := range a.AllowedNssaiList {
            for _, allowedSnssai := range allowedNssai.AllowedSnssaiList {
                if checkSnssaiInNssai(*allowedSnssai.AllowedSnssai, amfSetConfig.SupportedNssai) == true {
                    continue
                } else {
                    hitAllowedNssai = false
                    break
                }
            }
            if hitAllowedNssai == false {
                break
            }
        }

        if hitAllowedNssai == false {
            continue
        } else {
            // Add AMF Set to Authorized Network Slice Info
            if amfSetConfig.AmfList != nil && len(amfSetConfig.AmfList) != 0 {
                // List of candidate AMF(s) provided in configuration
                // TODO: Possibly querying the NRF
                a.CandidateAmfList = append(a.CandidateAmfList, amfSetConfig.AmfList...)
            } else {
                a.TargetAmfSet = amfSetConfig.AmfSetId
            }
            hitAmfSet = true
            break
        }
    }

    if hitAmfSet == false {
        flog.Warn("No AMF Set in configuration can serve the UE")
    }
}

// Use Subscribed S-NSSAI(s) which are marked as default S-NSSAI(s)
func useDefaultSubscribedSnssai(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo) {
    for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
        if subscribedSnssai.DefaultIndication == true {
            // Subscribed S-NSSAI is marked as default S-NSSAI

            var mappingOfSubscribedSnssai Snssai
            if p.HomePlmnId != nil && checkStandardSnssai(*subscribedSnssai.SubscribedSnssai) == false {
                // Find mapping of Subscribed S-NSSAI of UE's HPLMN to S-NSSAI in Serving PLMN from NSSF configuration
                mappingOfSnssai := getMappingOfPlmnFromConfig(*p.HomePlmnId)

                if mappingOfSnssai == nil {
                    flog.Warn("No S-NSSAI mapping of UE's HPLMN %+v in NSSF configuration", *p.HomePlmnId)
                    break
                }

                targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai,
                                                                  mappingOfSnssai)

                if found == false {
                    flog.Warn("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
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
            if p.HomePlmnId != nil {
                allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
            }

            // Default Access Type is set to 3GPP Access if no TAI is provided
            var accessType AccessType = IS_3_GPP_ACCESS
            if p.Tai != nil {
                accessType = getAccessTypeFromConfig(*p.Tai)
            }

            addAllowedSnssai(allowedSnssaiElement, accessType, a)

        }
    }

    addAmfInformation(a)
}

// Network slice selection for registration
// The function is executed when the IE, `slice-info-request-for-registration`, is provided in query parameters
func nsselectionForRegistration(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo, d *ProblemDetails) (status int) {
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

    if p.SliceInfoRequestForRegistration.RequestedNssai != nil {
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
                    if p.HomePlmnId != nil {
                        allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                        *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
                    }

                    // Default Access Type is set to 3GPP Access if no TAI is provided
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
        useDefaultSubscribedSnssai(p, a)
    }

    status = http.StatusOK
    return
}
