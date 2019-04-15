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
        flog.Warn("Invalid PLMN ID %+v provided in TAI", *tai.PlmnId)
        return false
    }

    for _, taConfig := range factory.NssfConfig.Configuration.TaList {
        if taConfig.Tac == tai.Tac {
            return true
        }
    }
    flog.Warn("No TA {Tac: %s} in NSSF configuration", tai.Tac)
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
            a.AllowedNssaiList[i].AllowedSnssaiList = append(a.AllowedNssaiList[i].AllowedSnssaiList, allowedSnssai)
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

// Use Subscribed S-NSSAI(s) which are marked as default S-NSSAI(s)
func useDefaultSubscribedSnssai(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo) {
    for _, subscribedSnssai := range p.SliceInfoRequestForRegistration.SubscribedNssai {
        if subscribedSnssai.DefaultIndication == true {
            // Subscribed S-NSSAI is marked as default S-NSSAI
            // Find mapping of Subscribed S-NSSAI of UE's HPLMN to S-NSSAI in Serving PLMN from NSSF configuration
            hitHomePlmn := false
            for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
                if *mappingFromPlmn.HomePlmnId == *p.HomePlmnId {
                    hitHomePlmn = true

                    targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai,
                                                                      mappingFromPlmn.MappingOfSnssai)

                    if found == false {
                        flog.Warn("No mapping of Subscribed S-NSSAI %+v in PLMN %+v in NSSF configuration",
                                  *subscribedSnssai.SubscribedSnssai,
                                  *p.HomePlmnId)
                        break
                    }

                    if checkSupportedSnssaiInTa(*targetMapping.ServingSnssai, p.Tai.Tac) == false {
                        continue
                    }

                    var allowedSnssaiElement AllowedSnssai
                    // TODO: Location configuration of NSI information list
                    allowedSnssaiElement.AllowedSnssai = new(Snssai)
                    *allowedSnssaiElement.AllowedSnssai = *targetMapping.ServingSnssai
                    if isRoamer == true {
                        allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                        *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
                    }

                    // TODO: Allowed NSSAI for different Access Type
                    var accessType AccessType = IS_3_GPP_ACCESS

                    addAllowedSnssai(allowedSnssaiElement, accessType, a)

                    break
                }
            }

            if hitHomePlmn == false {
                flog.Warn("No S-NSSAI mapping of UE's HPLMN %+v in NSSF configuration", *p.HomePlmnId)
                continue
            }
        }
    }
}

func nsselectionForRegistration(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo, d *ProblemDetails) (status int) {
    if isRoamer == true {
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
            if isRoamer == true && checkStandardSnssai(requestedSnssai) == false {
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
                    // TODO: Local configuration of NSI information list
                    allowedSnssaiElement.AllowedSnssai = new(Snssai)
                    *allowedSnssaiElement.AllowedSnssai = requestedSnssai
                    if isRoamer == true {
                        allowedSnssaiElement.MappedHomeSnssai = new(Snssai)
                        *allowedSnssaiElement.MappedHomeSnssai = *subscribedSnssai.SubscribedSnssai
                    }

                    // TODO: Allowed NSSAI for different Access Type
                    var accessType AccessType = IS_3_GPP_ACCESS

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

        if checkIfRequestAllowed == false {
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
