package ssnrgo

import (
	"encoding/binary"
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

const (
	timeFormat string = "2006-01-02 15:04:05"
)

func NewAnonymousNotification(receptorCode uint16, message string) *Notification {
	return NewNotificationWithTime(receptorCode, time.Now(), "Anonymous", message)
}

func NewNotification(receptorCode uint16, emitterName string, message string) *Notification {
	return NewNotificationWithTime(receptorCode, time.Now(), emitterName, message)
}

func NewNotificationWithTime(receptorCode uint16, timeStamp time.Time, emitterName string, message string) *Notification {
	r := new(Notification)
	r.oType = 78 // (N)otification's code
	r.rCode = receptorCode
	r.tStmp = timeStamp
	r.eName = emitterName
	r.message = message
	return r
}

func (n *Notification) GetCode() byte       { return n.oType }
func (n *Notification) GetReceptor() uint16 { return n.rCode }
func (n *Notification) GetTime() time.Time  { return n.tStmp }
func (n *Notification) GetEmiter() string   { return n.eName }
func (n *Notification) GetMessage() string  { return n.message }
func (n *Notification) GetSize() int {
	return 23 + len(n.message) + 1
}

func (n *Notification) Encode() []byte {
	r := make([]byte, n.GetSize())
	r[0] = n.oType
	binary.BigEndian.PutUint16(r[1:], n.rCode)
	binary.BigEndian.PutUint32(r[3:], uint32(n.tStmp.Unix()))
	copy(r[7:23], n.eName)
	copy(r[24:], n.message)
	return r
}

func DecodeNotification(array []byte) *Notification {
	r := new(Notification)
	r.oType = array[0]
	r.rCode = binary.BigEndian.Uint16(array[1:])
	r.tStmp = time.Unix(int64(binary.BigEndian.Uint32(array[3:])), 0)
	r.eName = string(array[7:23])
	r.message = string(array[24:])
	return r
}

func (n *Notification) String() string {
	return fmt.Sprintf("Notification sent at %s\nFrom: %s\nTo: %d\n%s",
		n.GetTime().Format(timeFormat),
		n.GetEmiter(),
		n.GetReceptor(),
		n.GetMessage())
}

func TestingThings() *Notification {
	return NewAnonymousNotification(313, "Things are good")
}
