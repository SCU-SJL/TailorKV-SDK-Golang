package tailorSDK

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	setex byte = iota
	setnx
	set
	get
	del
	unlink
	incr
	incrby
	ttl
	keys
	cnt
	save
	load
	cls
	exit
	quit
)

var tailorErrors = []error{
	errors.New(""),
	errors.New("syntax is wrong"),
	errors.New("item not found"),
	errors.New("item has already existed"),
	errors.New("no expiration cache of TailorKV save failed"),
	errors.New("expired cache of TailorKV save failed"),
	errors.New("load backup file failed"),
}

type Tailor struct {
	conn net.Conn
}

func (t *Tailor) Shutdown() error {
	return t.conn.Close()
}

func (t *Tailor) Set(key, val string) error {
	err := t.sendDatagram(set, key, val, "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Setex(key, val string, exp time.Duration) error {
	err := t.sendDatagram(setex, key, val, strconv.FormatInt(exp.Milliseconds(), 10))
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Setnx(key, val string) error {
	err := t.sendDatagram(setnx, key, val, "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Get(key string, maxSizeOfVal int) (string, error) {
	err := t.sendDatagram(get, key, "", "")
	if err != nil {
		return "", nil
	}

	err = t.readRespMsg()
	if err != nil {
		return "", err
	}

	val := make([]byte, maxSizeOfVal)
	n, err := t.read(val)
	if err != nil {
		return "", err
	}
	return string(val[:n]), nil
}

func (t *Tailor) Cnt() (int, error) {
	err := t.sendDatagram(cnt, "", "", "")
	if err != nil {
		return -1, err
	}
	if err = t.readRespMsg(); err != nil {
		return -1, err
	}
	count := make([]byte, 8)
	n, err := t.read(count)
	if err != nil {
		return -1, nil
	}
	cnt, err := strconv.Atoi(string(count[:n]))
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

func (t *Tailor) Del(key string) error {
	err := t.sendDatagram(del, key, "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Unlink(key string) error {
	err := t.sendDatagram(unlink, key, "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Ttl(key string) (time.Duration, error) {
	err := t.sendDatagram(ttl, key, "", "")
	if err != nil {
		return 0, err
	}

	err = t.readRespMsg()
	if err != nil {
		return 0, err
	}

	ttl := make([]byte, 128)
	n, err := t.read(ttl)
	if err != nil {
		return 0, err
	}

	return time.ParseDuration(string(ttl[:n]))
}

func (t *Tailor) sendDatagram(op byte, key, val, exp string) error {
	data := &datagram{
		Op:  op,
		Key: key,
		Val: val,
		Exp: exp,
	}
	dataJson, err := data.getJsonBytes()
	if err != nil {
		return err
	}
	_, err = t.conn.Write(dataJson)
	return err
}

func (t *Tailor) readRespMsg() error {
	errMsg := make([]byte, 128)
	n, err := t.read(errMsg)
	if err != nil {
		return err
	}
	if n == 1 && errMsg[0] != 0 {
		return tailorErrors[errMsg[0]]
	} else if n > 1 {
		return fmt.Errorf("TailorKV errMsg: %s\n", errMsg[:n])
	}
	return nil
}

func (t *Tailor) read(buf []byte) (int, error) {
	return t.conn.Read(buf)
}

func (t *Tailor) write(buf []byte) (int, error) {
	return t.conn.Write(buf)
}
