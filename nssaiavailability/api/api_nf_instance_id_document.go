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
	"free5gc-nssf/nssf_handler"
	"free5gc-nssf/nssf_handler/nssf_message"
)

func ApiNfInstanceIdDocumentDelete(c *gin.Context) {
	channelMsg := nssf_message.NewHttpChannelMessage()
	channelMsg.Event = nssf_message.EventNSSAIAvailabilityDelete
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

func ApiNfInstanceIdDocumentPatch(c *gin.Context) {
	channelMsg := nssf_message.NewHttpChannelMessage()
	channelMsg.Event = nssf_message.EventNSSAIAvailabilityPatch
	channelMsg.Context = c

	nssf_handler.SendMessage(channelMsg)
	rcvMsg := <-channelMsg.HttpChannel

	switch rcvMsg.Type {
	case nssf_message.HttpResponseMessageResponse:
		response := rcvMsg.Response.(AuthorizedNssaiAvailabilityInfo)
		c.JSON(http.StatusOK, response)
	case nssf_message.HttpResponseMessageProblemDetails:
		response := rcvMsg.Response.(ProblemDetails)
		c.JSON(int(response.Status), response)
	}
}

func ApiNfInstanceIdDocumentPut(c *gin.Context) {
	channelMsg := nssf_message.NewHttpChannelMessage()
	channelMsg.Event = nssf_message.EventNSSAIAvailabilityPut
	channelMsg.Context = c

	nssf_handler.SendMessage(channelMsg)
	rcvMsg := <-channelMsg.HttpChannel

	switch rcvMsg.Type {
	case nssf_message.HttpResponseMessageResponse:
		response := rcvMsg.Response.(AuthorizedNssaiAvailabilityInfo)
		c.JSON(http.StatusOK, response)
	case nssf_message.HttpResponseMessageProblemDetails:
		response := rcvMsg.Response.(ProblemDetails)
		c.JSON(int(response.Status), response)
	}
}
