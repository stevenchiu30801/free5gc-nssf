package nssf_handler

import (
	"time"

	"free5gc-nssf/nssaiavailability"
	"free5gc-nssf/nsselection"
	"free5gc-nssf/flog"
	"free5gc-nssf/nssf_handler/nssf_message"
)

const (
	MaxChannel int = 100000
)

var nssfChannel chan nssf_message.ChannelMessage

func init() {
	// init Pool
	nssfChannel = make(chan nssf_message.ChannelMessage, MaxChannel)
}

func SendMessage(msg nssf_message.ChannelMessage) {
	nssfChannel <- msg
}

func Handle() {
	for {
		select {
		case msg, ok := <-nssfChannel:
			if ok {
				switch msg.Event {
				case nssf_message.EventNSSelectionGet:
					nsselection.NSSelectionGet(msg.HttpChannel, msg.Context)
				case nssf_message.EventNSSAIAvailabilityPut:
					nssaiavailability.NSSAIAvailabilityPut(msg.HttpChannel, msg.Context)
				case nssf_message.EventNSSAIAvailabilityPatch:
					nssaiavailability.NSSAIAvailabilityPatch(msg.HttpChannel, msg.Context)
				case nssf_message.EventNSSAIAvailabilityDelete:
					nssaiavailability.NSSAIAvailabilityDelete(msg.HttpChannel, msg.Context)
				case nssf_message.EventNSSAIAvailabilityPost:
					nssaiavailability.NSSAIAvailabilityPost(msg.HttpChannel, msg.Context)
				case nssf_message.EventNSSAIAvailabilityUnsubscribe:
					nssaiavailability.NSSAIAvailabilityUnsubscribe(msg.HttpChannel, msg.Context)
				default:
					flog.Handler.Warnf("Event[%s] has not implemented", msg.Event)
				}
			} else {
				flog.Handler.Errorf("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
