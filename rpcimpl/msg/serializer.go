package msg

import (
	"encoding/json"
	"fmt"

	"tiny_rpc/log"
	"github.com/golang/protobuf/proto"
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
		defSerializer = JsonSerializer{}
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

func (r JsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (r JsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type PBSerializer struct {
}

func (r PBSerializer) Marshal(v interface{}) ([]byte, error) {
	m, ok := v.(proto.Message)
	if ok {
		return nil, fmt.Errorf("not proto msg")
	}
	return proto.Marshal(m)
}

func (r PBSerializer) Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if ok {
		return fmt.Errorf("not proto msg")
	}
	return proto.Unmarshal(data, m)
}
