/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf

import (
    "net/http"

    . "../model"
)

func findSnssaiFromMapping(mappings []MappingOfSnssai, s Snssai) (MappingOfSnssai, bool) {
    for _, m := range mappings {
        if *m.ServingSnssai == s {
            return m, true
        }
    }
    return MappingOfSnssai{}, false
}

func nsselectionForRegistration(p NsselectionQueryParameter, a *AuthorizedNetworkSliceInfo, d *ProblemDetails) int {
    if p.SliceInfoRequestForRegistration.RequestedNssai != nil {
        for _, r := range p.SliceInfoRequestForRegistration.RequestedNssai {
            targetMapping, found := findSnssaiFromMapping(p.SliceInfoRequestForRegistration.MappingOfNssai, r)

            if found == false {
                // TODO: Search for local configuration if there is no provided mapping from UE
                continue
            }

            for _, s := range p.SliceInfoRequestForRegistration.SubscribedNssai {
                if *targetMapping.HomeSnssai == *s.SubscribedSnssai {
                    var allowedNssai AllowedNssai
                    var allowedSnssaiList AllowedSnssai
                    var allowedSnssai, mappedHomeSnssai Snssai = r, *s.SubscribedSnssai

                    // TODO: Location configuration of NSI information list
                    allowedSnssaiList.AllowedSnssai = &allowedSnssai
                    allowedSnssaiList.MappedHomeSnssai = &mappedHomeSnssai

                    allowedNssai.AllowedSnssaiList = append(allowedNssai.AllowedSnssaiList, allowedSnssaiList)
                    // TODO: Allowed NSSAI for different Access Type
                    var accessType AccessType = IS_3_GPP_ACCESS
                    allowedNssai.AccessType = &accessType

                    a.AllowedNssaiList = append(a.AllowedNssaiList, allowedNssai)
                }
            }
            // TODO: Operation if no mapping between Requested Snssai and Snssai of UE's HPLMN
        }
    } else {
        // for _, s := p.SliceInfoForRegistration.SubscribedNssai {
        //
        // }
    }

    return http.StatusOK
}
