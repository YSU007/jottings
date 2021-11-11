package msg

import (
	"encoding/json"
	"fmt"

	"Jottings/tiny_rpc/log"
)

type SerializerType byte

const (
	SerializerJson SerializerType = iota
	SerializerPB
)

var defSerializer Serializer

func SetSerializer(t SerializerType) {
	switch t {
	case SerializerJson:
		defSerializer = JsonSerializer{}
	case SerializerPB:

	default:
		log.Error("Serializer not find,type %v", t)
	}

}

func Marshal(v interface{}) ([]byte, error) {
	if defSerializer == nil {
		return nil, fmt.Errorf("defSerializer nil")
	}
	return defSerializer.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if defSerializer == nil {
		return fmt.Errorf("defSerializer nil")
	}
	return defSerializer.Unmarshal(data, v)
}

type JsonSerializer struct {
}

func (j JsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j JsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
