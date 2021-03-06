package tailorSDK

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
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
	mu   sync.Mutex
}

func (t *Tailor) Set(key, val string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(set, key, val, "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Setex(key, val string, exp time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(setex, key, val, strconv.FormatInt(exp.Milliseconds(), 10))
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Setnx(key, val string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(setnx, key, val, "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Get(key string, maxSizeOfVal int) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
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
	t.mu.Lock()
	defer t.mu.Unlock()
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
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(del, key, "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Unlink(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(unlink, key, "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Ttl(key string) (time.Duration, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
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

func (t *Tailor) Incr(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(incr, key, "", "")
	if err != nil {
		return err
	}

	return t.readRespMsg()
}

func (t *Tailor) Incrby(key string, addition int) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(incrby, key, strconv.Itoa(addition), "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Cls() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(cls, "", "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

func (t *Tailor) Save(filename string) error {
	err := t.sendDatagram(save, filename, "", "")
	if err != nil {
		return err
	}
	err = t.readRespMsg()
	if err != nil {
		return errors.New("no expiration part of tailor kv failed to save")
	}
	err = t.readRespMsg()
	if err != nil {
		return errors.New("expiry part of tailor kv failed to save")
	}
	return nil
}

func (t *Tailor) Load(filename string) error {
	err := t.sendDatagram(load, filename, "", "")
	if err != nil {
		return err
	}
	return t.readRespMsg()
}

type keysDatagram struct {
	Ks []string `json:"keys"`
}

func getKeys(data []byte) ([]string, error) {
	var ks keysDatagram
	err := json.Unmarshal(data, &ks)
	if err != nil {
		return nil, err
	}
	return ks.Ks, nil
}

func (t *Tailor) Keys(expr string, maxSize int) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.sendDatagram(keys, expr, "", "")
	if err != nil {
		return nil, err
	}
	err = t.readRespMsg()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, maxSize)
	n, err := t.read(buf)
	if err != nil {
		return nil, err
	}
	return getKeys(buf[:n])
}

func (t *Tailor) Shutdown() error {
	t.mu.Lock()
	defer func() {
		t.mu.Unlock()
		t.conn.Close()
	}()
	return t.sendDatagram(exit, "", "", "")
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
