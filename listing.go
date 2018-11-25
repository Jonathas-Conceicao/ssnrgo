package ssnrgo

import (
	"encoding/binary"
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

func (l *Listing) Encode() []byte {
	r := make([]byte, l.GetSize())
	if l.users != nil {
		l.amount = l.users.PutTo(r[5:], l.offset, l.amount)
	}
	r[0] = l.oType
	binary.BigEndian.PutUint16(r[1:], l.amount)
	binary.BigEndian.PutUint16(r[3:], l.offset)
	return r
}

func DecodeListing(data []byte) *Listing {
	r := new(Listing)
	r.oType = data[0]
	r.amount = binary.BigEndian.Uint16(data[1:])
	r.offset = binary.BigEndian.Uint16(data[3:])
	r.users = nil
	return r
}

func DecodeListingReceived(data []byte) *Listing {
	r := DecodeListing(data)
	usrs := new(UserTable)
	data = data[5:]
	for i := 0; i < int(r.amount)*UserSize; i += UserSize {
		usrs.Add(
			binary.BigEndian.Uint16(data[i:]),
			User{string(data[i+2 : i+UserSize]), nil})
	}
	r.SetUsers(usrs)
	return r
}

func (l *Listing) String() string {
	return l.users.String()
}
