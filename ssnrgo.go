package ssnrgo

import (
	"fmt"
	"time"
)

type Notification struct {
	oType   uint8
	rCode   uint16
	tStmp   time.Time
	eName   string
	message string
}

func NewNotification(
	receptorCode uint16,
	timeStamp time.Time,
	emitterName string,
	message string,
) *Notification {
	r := new(Notification)
	r.oType = 78 // (N)otification's code
	r.rCode = receptorCode
	r.tStmp = timeStamp
	r.eName = fmt.Sprintf("%16s", emitterName)
	r.message = message
	return r
}

func (n *Notification) getMessage() string {
	return n.message
}
func (n *Notification) getTime() time.Time {
	return n.tStmp
}

func TestingThings() string {
	message := NewNotification(0, time.Now(), "", "Things were good")
	return message.getMessage() + " at " +
		message.getTime().Format("2006-01-02 15:04:05")
}
