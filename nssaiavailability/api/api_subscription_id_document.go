/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "free5gc-nssf/model"
	"free5gc-nssf/flog"
	"free5gc-nssf/nssf_handler"
	"free5gc-nssf/nssf_handler/nssf_message"
)

func ApiSubscriptionIdDocument(c *gin.Context) {
	// Due to conflict of route matching, 'subscriptions' in the route is replaced with the existing wildcard ':nfId'
	nfId := c.Param("nfId")
	if nfId != "subscriptions" {
		c.JSON(http.StatusNotFound, gin.H{})
		flog.Nssaiavailability.Infof("404 Not Found")
		return
	}

	channelMsg := nssf_message.NewHttpChannelMessage()
	channelMsg.Event = nssf_message.EventNSSAIAvailabilityUnsubscribe
	channelMsg.Context = c

	nssf_handler.SendMessage(channelMsg)
	rcvMsg := <-channelMsg.HttpChannel

	switch rcvMsg.Type {
	case nssf_message.HttpResponseMessageResponse:
		c.JSON(http.StatusNoContent, gin.H{})
	case nssf_message.HttpResponseMessageProblemDetails:
		response := rcvMsg.Response.(ProblemDetails)
		c.JSON(int(response.Status), response)
	}
}
