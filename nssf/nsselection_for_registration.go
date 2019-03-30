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
            for _, mappingSet := range factory.NssfConfig.Configuration.MappingSet {
                if *mappingSet.HomePlmnId == *p.HomePlmnId {
                    hitHomePlmn = true

                    targetMapping, found := findMappingWithHomeSnssai(*subscribedSnssai.SubscribedSnssai,
                                                                      mappingSet.MappingOfSnssai)

                    if found == false {
                        flog.Warn("No mapping of Subscribed S-NSSAI %+v in PLMN %+v",
                                  *subscribedSnssai.SubscribedSnssai,
                                  *p.HomePlmnId)
                        break
                    }

                    var allowedSnssaiElement AllowedSnssai
                    var allowedSnssai Snssai = *targetMapping.ServingSnssai
                    var mappedHomeSnssai Snssai = *subscribedSnssai.SubscribedSnssai
                    // TODO: Location configuration of NSI information list
                    allowedSnssaiElement.AllowedSnssai = &allowedSnssai
                    allowedSnssaiElement.MappedHomeSnssai = &mappedHomeSnssai

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

func nsselectionForRegistration(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo, d *ProblemDetails) int {
    if p.SliceInfoRequestForRegistration.RequestedNssai != nil {
        // Requested NSSAI is provided
        // Verify which S-NSSAI(s) in the Requested NSSAI are permitted based on comparing the Subscribed S-NSSAI(s)

        // Check if any Requested S-NSSAIs is present in Subscribed S-NSSAIs
        checkIfRequestAllowed := false

        for _, requestedSnssai := range p.SliceInfoRequestForRegistration.RequestedNssai {
            isStandardSnssai := checkStandardSnssai(requestedSnssai)
            var mappingOfRequestedSnssai Snssai
            if isStandardSnssai == false {
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
                    var allowedSnssai Snssai = requestedSnssai
                    var mappedHomeSnssai Snssai = *subscribedSnssai.SubscribedSnssai
                    // TODO: Location configuration of NSI information list
                    allowedSnssaiElement.AllowedSnssai = &allowedSnssai
                    allowedSnssaiElement.MappedHomeSnssai = &mappedHomeSnssai

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

    return http.StatusOK
}
