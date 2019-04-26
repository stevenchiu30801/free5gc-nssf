/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
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

// Check whether the given S-NSSAI is supported or not in PLMN
func checkSupportedSnssaiInPlmn(snssai Snssai) bool {
    for _, supportedSnssai := range factory.NssfConfig.Configuration.SupportedNssaiInPlmn {
        if snssai == supportedSnssai {
            return true
        }
    }
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
                // The API URI of the NRF may be included if target AMF Set is included
                a.NrfAmfSet = amfSetConfig.NrfAmfSet
            }
            hitAmfSet = true
            break
        }
    }

    if hitAmfSet == false {
        flog.Warn("No AMF Set in configuration can serve the UE")
    }
}
