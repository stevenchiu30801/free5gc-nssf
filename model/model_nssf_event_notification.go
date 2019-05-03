/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

type NssfEventNotification struct {

	SubscriptionId string `json:"subscriptionId"`

	AuthorizedNssaiAvailabilityData []AuthorizedNssaiAvailabilityData `json:"authorizedNssaiAvailabilityData"`
}
