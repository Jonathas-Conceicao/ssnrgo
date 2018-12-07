package ssnrgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Notification struct {
	oType   uint8
	mSize   uint64
	rCode   uint16
	tStmp   time.Time
	eName   string
	message string
}

const (
	NotificationCode uint8  = 78 // (N)otification's code
	timeFormat       string = "2006-01-02 15:04:05"

	NotificationHeaderSize int = 31
)

func NewNotification(receptorCode uint16, emitterName string, message string) *Notification {
	if emitterName != "" {
		return NewNotificationWithTime(receptorCode, time.Now(), emitterName, message)
	} else {
		return NewNotificationWithTime(receptorCode, time.Now(), "Anonymous", message)
	}
}

func NewNotificationWithTime(receptorCode uint16, timeStamp time.Time, emitterName string, message string) *Notification {
	r := new(Notification)
	r.oType = NotificationCode
	r.rCode = receptorCode
	r.tStmp = timeStamp
	r.eName = emitterName
	r.message = message
	return r
}

func (n *Notification) GetCode() byte       { return n.oType }
func (n *Notification) GetReceptor() uint16 { return n.rCode }
func (n *Notification) GetTime() time.Time  { return n.tStmp }
func (n *Notification) GetEmitter() string  { return n.eName }
func (n *Notification) GetMessage() string  { return n.message }
func (n *Notification) GetSize() int        { return len(n.message) + 1 }
func (n *Notification) GetTimeString() string {
	return n.tStmp.Format(timeFormat)
}

func (n *Notification) Encode() []byte {
	r := make([]byte, n.GetSize()+NotificationHeaderSize)
	r[0] = n.oType
	binary.BigEndian.PutUint64(r[1:], uint64(n.GetSize()))
	binary.BigEndian.PutUint16(r[9:], n.rCode)
	binary.BigEndian.PutUint32(r[11:], uint32(n.tStmp.Unix()))
	copy(r[15:31], n.eName)
	copy(r[32:], n.message)
	return r
}

func DecodeNotification(array []byte) (*Notification, error) {
	if array[0] != NotificationCode {
		return nil, errors.New("Invalid Notification code: " +
			strconv.FormatInt(int64(array[0]), 10))
	}
	r := new(Notification)
	r.oType = array[0]
	r.rCode = binary.BigEndian.Uint16(array[9:])
	r.tStmp = time.Unix(int64(binary.BigEndian.Uint32(array[11:])), 0)
	r.eName = string(bytes.Trim(array[15:31], "\x00"))
	r.message = string(array[32:])
	return r, nil
}

func ReadNotification(rd *bufio.Reader) (
	[]byte, *Notification, error) {
	slice, err := rd.Peek(9)
	if err != nil {
		return slice, nil, err
	}
	size := binary.BigEndian.Uint64(slice[1:]) + uint64(NotificationHeaderSize)
	data := make([]byte, size)
	_, err = rd.Read(data)
	if err != nil {
		return data, nil, err
	}
	ntf, err := DecodeNotification(data)
	return data, ntf, err
}

func (n *Notification) String() string {
	return fmt.Sprintf("%s -- From: \"%s\"\n%s",
		n.GetTime().Format(timeFormat),
		n.GetEmitter(),
		n.GetMessage())
}
