package tailorSDK

import (
	"errors"
	"fmt"
	"net"
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
