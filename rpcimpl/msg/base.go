package msg

import (
	"encoding/binary"
	"io"

	"tiny_rpc/log"
)

type Head struct {
	Len   uint32
	MType MType
}

type ModeBase struct {
	Head
	Mode uint32
	Data []byte
}

func (r *ModeBase) MsgType() MType {
	return r.MType
}

func (r *ModeBase) FillIn(mode uint32, data []byte) {
	log.Error("msg type illegal %v mode %v", MTypeIllegal, mode)
}

func (r *ModeBase) GetMode() uint32 {
	return r.Mode
}

func (r *ModeBase) GetData() []byte {
	return r.Data
}

func (r *ModeBase) Encode(writer io.Writer) error {
	var streamSlice = make([]byte, r.Len+4+1+4)
	binary.BigEndian.PutUint32(streamSlice[0:4], r.Len)
	streamSlice[4] = byte(r.MType)
	binary.BigEndian.PutUint32(streamSlice[5:9], r.Mode)
	copy(streamSlice[9:], r.Data)
	_, err := writer.Write(streamSlice)
	return err
}

func (r *ModeBase) Decode(reader io.Reader) error {
	var slice4byte = make([]byte, 4)
	_, err := io.ReadFull(reader, slice4byte)
	if err != nil {
		return err
	}
	var l = binary.BigEndian.Uint32(slice4byte)

	var t = make([]byte, 1)
	_, err = io.ReadFull(reader, t)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(reader, slice4byte)
	if err != nil {
		return err
	}
	var m = binary.BigEndian.Uint32(slice4byte)

	var payload = make([]byte, l)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return err
	}

	r.Len = l
	r.MType = MType(t[0])
	r.Mode = m
	r.Data = payload
	return nil
}

type CodeBase struct {
	Head
	Code uint32
	Data []byte
}

func (r *CodeBase) MsgType() MType {
	return r.MType
}

func (r *CodeBase) FillIn(code uint32, data []byte) {
	log.Error("msg type illegal")
}

func (r *CodeBase) GetCode() uint32 {
	return r.Code
}

func (r *CodeBase) GetData() []byte {
	return r.Data
}

func (r *CodeBase) Encode(writer io.Writer) error {
	var streamSlice = make([]byte, r.Len+4+1+4)
	binary.BigEndian.PutUint32(streamSlice[0:4], r.Len)
	streamSlice[4] = byte(r.MType)
	binary.BigEndian.PutUint32(streamSlice[5:9], r.Code)
	copy(streamSlice[9:], r.Data)
	_, err := writer.Write(streamSlice)
	return err
}

func (r *CodeBase) Decode(reader io.Reader) error {
	var slice4byte = make([]byte, 4)
	_, err := io.ReadFull(reader, slice4byte)
	if err != nil {
		return err
	}
	var l = binary.BigEndian.Uint32(slice4byte)

	var t = make([]byte, 1)
	_, err = io.ReadFull(reader, t)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(reader, slice4byte)
	if err != nil {
		return err
	}
	var c = binary.BigEndian.Uint32(slice4byte)

	var payload = make([]byte, l)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return err
	}

	r.Len = l
	r.MType = MType(t[0])
	r.Code = c
	r.Data = payload
	return nil
}

// NotifyBase PushBase ----------------------------------------------------------------------------------------------------
type NotifyBase struct {
	ModeBase
}

func (r *NotifyBase) FillIn(mode uint32, data []byte) {
	r.Len = uint32(len(data))
	r.MType = MTypeNotice
	r.Mode = mode
	r.Data = data
}

type PushBase struct {
	ModeBase
}

func (r *PushBase) FillIn(mode uint32, data []byte) {
	r.Len = uint32(len(data))
	r.MType = MTypePush
	r.Mode = mode
	r.Data = data
}

// RequestBase ResponseBase ----------------------------------------------------------------------------------------------------
type RequestBase struct {
	ModeBase
}

func (r *RequestBase) FillIn(mode uint32, data []byte) {
	r.Len = uint32(len(data))
	r.MType = MTypeRpc
	r.Mode = mode
	r.Data = data
}

type ResponseBase struct {
	CodeBase
}

func (r *ResponseBase) FillIn(code uint32, data []byte) {
	r.Len = uint32(len(data))
	r.MType = MTypeRpc
	r.Code = code
	r.Data = data
}
