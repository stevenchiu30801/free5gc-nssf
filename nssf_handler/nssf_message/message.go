package nssf_message

import (
	"github.com/gin-gonic/gin"
)

type HttpResponseMessageType string

const (
	HttpResponseMessageResponse       HttpResponseMessageType = "Response"
	HttpResponseMessageProblemDetails HttpResponseMessageType = "Problem Details"
)

type ChannelMessage struct {
	Event       int
	HttpChannel chan HttpResponseMessage // Channel of returned Http response
	Context     *gin.Context             // Request context of Gin server
}

func NewHttpChannelMessage() ChannelMessage {
	msg := ChannelMessage{}
	msg.HttpChannel = make(chan HttpResponseMessage)
	return msg
}

type HttpResponseMessage struct {
	Type     HttpResponseMessageType
	Response interface{}
}

// Send HTTP Response to HTTP handler thread through HTTP channel
func SendHttpResponseMessage(channel chan HttpResponseMessage, responseType HttpResponseMessageType, args interface{}) {
	responseMsg := HttpResponseMessage{}
	responseMsg.Type = responseType
	responseMsg.Response = args

	channel <- responseMsg
}
