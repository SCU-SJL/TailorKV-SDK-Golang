package tailorSDK

import "encoding/json"

type datagram struct {
	Op  byte   `json:"op"`
	Key string `json:"key"`
	Val string `json:"val,omitempty"`
	Exp string `json:"exp,omitempty"`
}

func (d *datagram) getJsonBytes() ([]byte, error) {
	return json.Marshal(*d)
}
