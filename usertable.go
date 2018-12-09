package ssnrgo

import (
	"encoding/binary"
	"net"
	"strconv"
	"sync"
	"time"
)

type User struct {
	Name string
	Addr net.Conn
}

type UserTable struct {
	set sync.Map
}

const (
	PingCode byte = 80 //(P)ing code

	UserSize = 18
)

func (t *UserTable) Length() int {
	length := 0
	t.set.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func (t *UserTable) GetSize() int {
	l := t.Length()
	return l * (2 + 16)
}

func (t *UserTable) Add(idx uint16, v User) (uint16, error) {
	ctrl := true
	for ctrl {
		_, ctrl = t.set.LoadOrStore(idx, v)
		idx += 1
	}
	return idx - 1, nil
}

func (t *UserTable) Get(idx uint16) *User {
	usr := new(User)
	usrI, ok := t.set.Load(idx)
	if ok {
		*usr = usrI.(User)
		return usr
	}
	return nil
}

func (t *UserTable) PutUsers(data []byte, offset, amount uint16) uint16 {
	if int(offset) >= t.Length() {
		return 0
	}
	n, i := uint16(0), 0
	t.set.Range(func(k, v interface{}) bool {
		if n < offset { // Skeep until offset
			n += 1
			return true
		}
		binary.BigEndian.PutUint16(data[i:], k.(uint16))
		copy(data[i+2:i+UserSize], v.(User).Name)
		i += UserSize
		n += 1
		return n < amount
	})
	return n - offset
}

func (t *UserTable) ClearDisconnects() int {
	var n int
	t.set.Range(func(k, v interface{}) bool {
		cn := v.(User).Addr
		if cn != nil {
			cn.SetReadDeadline(time.Now())
			cn.Write([]byte{PingCode})
			_, err := cn.Write([]byte{PingCode})
			if err != nil {
				cn.Close()
				t.set.Delete(k)
				n++
			} else {
				cn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			}
		}
		return true
	})
	return n
}

func (t *UserTable) String() string {
	users := ""
	t.set.Range(func(k, v interface{}) bool {
		users += strconv.Itoa(int(k.(uint16))) + ":" + v.(User).Name + "\n"
		return true
	})
	return users
}
