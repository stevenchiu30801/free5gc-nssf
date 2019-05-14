/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "encoding/json"
    "reflect"

    factory "../factory"
    flog "../flog"
    . "../model"
)

// Check if a slice contains an element
func Contain(target interface{}, slice interface{}) bool {
    arr := reflect.ValueOf(slice)
    if arr.Kind() == reflect.Slice {
        for i := 0; i < arr.Len(); i++ {
            if reflect.DeepEqual(arr.Index(i).Interface(), target) {
                return true
            }
        }
    }
    return false
}

// Check whether UE's Home PLMN is configured/supported
func checkSupportedHplmn(homePlmnId PlmnId) bool {
    for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
        if *mappingFromPlmn.HomePlmnId == homePlmnId {
            return true
        }
    }
    flog.Nsselection.Warnf("No Home PLMN %+v in NSSF configuration", homePlmnId)
    return false
}

// Check whether UE's current TA is configured/supported
func checkSupportedTa(tai Tai) bool {
    // TODO: Check supported TA under given AMF
    if *tai.PlmnId != *factory.NssfConfig.Info.ServingPlmnId {
        flog.Nsselection.Infof("Invalid PLMN ID %+v provided in TAI", *tai.PlmnId)
        return false
    }

    for _, taConfig := range factory.NssfConfig.Configuration.TaList {
        if reflect.DeepEqual(*taConfig.Tai, tai) == true {
            return true
        }
    }
    e, _ := json.Marshal(&tai)
    flog.Nsselection.Warnf("No TA %s in NSSF configuration", e)
    return false
}

// Check whether the given S-NSSAI is supported or not in PLMN
func checkSupportedSnssaiInPlmn(snssai Snssai) bool {
    if checkStandardSnssai(snssai) == true {
        return true
    }

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

// Check whether S-NSSAI is in SupportedNssaiAvailabilityData under the given TAI
func checkSupportedNssaiAvailabilityData(snssai Snssai, tai Tai, s []SupportedNssaiAvailabilityData) bool {
    for _, supportedNssaiAvailabilityData := range s {
        if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) == true &&
           Contain(snssai, supportedNssaiAvailabilityData.SupportedSnssaiList) == true {
            return true
        }
    }
    return false
}

// Check whether S-NSSAI is supported or not in UE's current TA
func checkSupportedSnssaiInTa(snssai Snssai, nfId string, tai Tai) bool {
    for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
        if amfSetConfig.AmfList != nil && len(amfSetConfig.AmfList) != 0 && Contain(nfId, amfSetConfig.AmfList) {
            return checkSupportedNssaiAvailabilityData(snssai, tai, amfSetConfig.SupportedNssaiAvailabilityData)
        }
    }

    for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
        if amfConfig.NfId == nfId {
            return checkSupportedNssaiAvailabilityData(snssai, tai, amfConfig.SupportedNssaiAvailabilityData)
        }
    }

    flog.Nsselection.Warnf("No AMF %s in NSSF configuration", nfId)
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
        if reflect.DeepEqual(*taConfig.Tai, tai) == true {
            return *taConfig.AccessType
        }
    }
    e, _ := json.Marshal(&tai)
    flog.Nsselection.Warnf("No TA %s in NSSF configuration", e)
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
                flog.Nsselection.Infof("Unable to add a new Allowed S-NSSAI since already eight S-NSSAIs in Allowed NSSAI")
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

// // Add AMF information to Authorized Network Slice Info
// func addAmfInformation(a *AuthorizedNetworkSliceInfo) {
//     if a.AllowedNssaiList == nil || len(a.AllowedNssaiList) == 0 {
//         return
//     }
//
//     // Check if any AMF can serve the UE
//     // That is, whether NSSAI of all Allowed S-NSSAIs is a subset of NSSAI supported by AMF
//     // Simply use the first applicable AMF set
//     // TODO: Policies of AMF selection (e.g. load balance between AMF instances)
//     hitAmfSet := false
//     for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
//         hitAllowedNssai := true
//         for _, allowedNssai := range a.AllowedNssaiList {
//             for _, allowedSnssai := range allowedNssai.AllowedSnssaiList {
//                 if Contain(*allowedSnssai.AllowedSnssai, amfSetConfig.SupportedNssai) == true {
//                     continue
//                 } else {
//                     hitAllowedNssai = false
//                     break
//                 }
//             }
//             if hitAllowedNssai == false {
//                 break
//             }
//         }
//
//         if hitAllowedNssai == false {
//             continue
//         } else {
//             // Add AMF Set to Authorized Network Slice Info
//             if amfSetConfig.AmfList != nil && len(amfSetConfig.AmfList) != 0 {
//                 // List of candidate AMF(s) provided in configuration
//                 // TODO: Possibly querying the NRF
//                 a.CandidateAmfList = append(a.CandidateAmfList, amfSetConfig.AmfList...)
//             } else {
//                 a.TargetAmfSet = amfSetConfig.AmfSetId
//                 // The API URI of the NRF may be included if target AMF Set is included
//                 a.NrfAmfSet = amfSetConfig.NrfAmfSet
//             }
//             hitAmfSet = true
//             break
//         }
//     }
//
//     if hitAmfSet == false {
//         flog.Nsselection.Warnf("No AMF Set in configuration can serve the UE")
//     }
// }
