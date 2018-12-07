package ssnrgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

type Listing struct {
	oType  uint8
	amount uint16
	offset uint16
	users  *UserTable
}

const (
	ListingCode uint8 = 76 // (L)isting's code
)

func NewListingRequestAll() *Listing {
	return NewListingRequest(^uint16(0), 0)
}

func NewListingRequest(amount, offset uint16) *Listing {
	r := new(Listing)
	r.oType = ListingCode
	r.amount = amount
	r.offset = offset
	r.users = nil
	return r
}

func (l *Listing) GetCode() byte        { return l.oType }
func (l *Listing) GetUsers() *UserTable { return l.users }
func (l *Listing) GetOffset() uint16    { return l.offset }
func (l *Listing) GetAmount() uint16    { return l.amount }
func (l *Listing) GetSize() int {
	size := 5 // 1 + 2 + 2
	if l.users != nil {
		size += l.users.Length() * UserSize
	}
	return size
}

func (l *Listing) SetUsers(users *UserTable) { l.users = users }

func (l *Listing) ReadUsers(data []byte) {
	usrs := new(UserTable)
	for i := 0; i < int(l.amount)*UserSize; i += UserSize {
		usrs.Add(
			binary.BigEndian.Uint16(data[i:]),
			User{string(bytes.Trim(data[i+2:i+UserSize], "\x00")), nil})
	}
	l.users = usrs
}

func (l *Listing) Encode() []byte {
	r := make([]byte, l.GetSize())
	if l.users != nil {
		l.amount = l.users.PutUsers(r[5:], l.offset, l.amount)
	}
	r[0] = l.oType
	binary.BigEndian.PutUint16(r[1:], l.amount)
	binary.BigEndian.PutUint16(r[3:], l.offset)
	return r
}

func DecodeListing(data []byte) (*Listing, error) {
	if data[0] != ListingCode {
		return nil, errors.New("Invalid Listing code: " +
			strconv.FormatInt(int64(data[0]), 10))
	}
	r := new(Listing)
	r.oType = data[0]
	r.amount = binary.BigEndian.Uint16(data[1:])
	r.offset = binary.BigEndian.Uint16(data[3:])
	r.users = nil
	return r, nil
}

func ReadListing(rd *bufio.Reader, request bool) (
	[]byte, *Listing, error) {
	slice, err := rd.Peek(3)
	if err != nil {
		return slice, nil, err
	}

	size := 5
	if !request {
		size += int(binary.BigEndian.Uint16(slice[1:])) * UserSize
	}
	data := make([]byte, size)
	_, err = rd.Read(data)
	if err != nil {
		return data, nil, err
	}

	lst, err := DecodeListing(data)
	if !request {
		lst.ReadUsers(data[5:])
	}
	return data, lst, err
}

func (l *Listing) String() string {
	return l.users.String()
}
