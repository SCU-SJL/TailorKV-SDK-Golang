package tailorSDK

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
	"net"
	"sync"
)

func Connect(ipAddr, port, password, AESKey string) (*Tailor, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipAddr+":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	authErr := auth(conn, password, AESKey)
	if authErr != nil {
		return nil, authErr
	}
	return &Tailor{
		conn: conn,
		mu:   sync.Mutex{},
	}, nil
}

func auth(conn net.Conn, password, AESKey string) error {
	need, err := needAuth(conn)
	if err != nil {
		return err
	}
	if need {
		password = aesEncrypt(password, AESKey)
		_, err = conn.Write([]byte(password))
		if err != nil {
			return err
		}
		resp := make([]byte, 1)
		_, err = conn.Read(resp)
		if err != nil {
			return err
		}
		if resp[0] != 0 {
			return errors.New("wrong password")
		}
	}
	return nil
}

func needAuth(conn net.Conn) (bool, error) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return true, err
	}
	return buf[0] == 1, nil
}

func aesEncrypt(orig string, key string) string {
	origData := []byte(orig)
	k := []byte(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		log.Fatal(err)
	}

	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)

	return base64.StdEncoding.EncodeToString(encrypted)
}

func pKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
