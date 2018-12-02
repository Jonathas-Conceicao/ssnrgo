package ssnrgo

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
)

type Register struct {
	oType uint8
	rCode uint16
	value uint8
	name  string
}

const (
	RegisterCode   uint8 = 82 // (R)egister's code
	DisconnectCode uint8 = 68 // (D)isconnect code

	RegisterMessageSize int = 20

	ConnAccepted  uint8 = 0
	ConnNewAddres uint8 = 1
	RefServerFull uint8 = 2
	RefBlackList  uint8 = 4
	RefUnknowEror uint8 = 128
)

func NewRegister(receptorCode uint16, name string) *Register {
	r := new(Register)
	r.oType = RegisterCode
	r.rCode = receptorCode
	r.name = name
	return r
}

func NewDisconnect() *Register {
	r := new(Register)
	r.oType = DisconnectCode
	return r
}

func (n *Register) GetReceptor() uint16  { return n.rCode }
func (n *Register) GetReturn() uint8     { return n.value }
func (n *Register) GetName() string      { return n.name }
func (n *Register) SetReturn(v uint8)    { n.value = v }
func (n *Register) SetReceptor(v uint16) { n.rCode = v }

func (n *Register) Encode() []byte {
	r := make([]byte, RegisterMessageSize)
	r[0] = n.oType
	binary.BigEndian.PutUint16(r[1:], n.rCode)
	r[3] = n.value
	copy(r[4:20], n.name)
	return r
}

func DecodeRegister(array []byte) (*Register, error) {
	if array[0] != RegisterCode {
		return nil, errors.New("Invalid Registration Code")
	}
	r := new(Register)
	r.oType = array[0]
	r.rCode = binary.BigEndian.Uint16(array[1:])
	r.value = array[3]
	r.name = string(array[4:])
	return r, nil
}

func ReadRegister(rd *bufio.Reader) (
	[]byte, *Register, error) {
	data := make([]byte, RegisterMessageSize)
	_, err := rd.Read(data)
	if err != nil {
		return data, nil, err
	}
	rgt, err := DecodeRegister(data)
	return data, rgt, err
}

func (n *Register) String() string {
	return fmt.Sprintf("Register at id %d", n.GetReceptor())
}
