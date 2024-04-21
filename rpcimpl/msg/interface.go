package msg

import "io"

type MType byte

const (
	MTypeIllegal MType = iota
	MTypeRpc
	MTypeNotice
	MTypePush
)

// ModeMsg CodeMsg ----------------------------------------------------------------------------------------------------
type ModeMsg interface {
	MsgType() MType
	FillIn(mode uint32, data []byte)
	GetMode() uint32
	GetData() []byte
	Encode(r io.Writer) error
	Decode(r io.Reader) error
}

type CodeMsg interface {
	MsgType() MType
	FillIn(code uint32, data []byte)
	GetCode() uint32
	GetData() []byte
	Encode(r io.Writer) error
	Decode(r io.Reader) error
}

// Serializer ----------------------------------------------------------------------------------------------------
type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}
