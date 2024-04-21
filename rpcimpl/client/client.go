package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"tiny_rpc/log"
	"tiny_rpc/msg"
)

type Client struct {
	conn net.Conn
}

func NewClient(network, address string) *Client {
	var conn, err = net.DialTimeout(network, address, time.Second)
	if err != nil {
		log.Error("NewClient err %v", err)
	}
	return &Client{
		conn: conn,
	}
}

func (c *Client) Call(mode uint32, req interface{}, rsp interface{}) (uint32, error) {
	var baseReq = new(msg.RequestBase)
	var data, _ = msg.Marshal(req)
	baseReq.FillIn(mode, data)
	var streamSlice = make([]byte, baseReq.Len+4+1+4)
	binary.BigEndian.PutUint32(streamSlice[0:4], baseReq.Len)
	streamSlice[4] = byte(baseReq.MType)
	binary.BigEndian.PutUint32(streamSlice[5:9], baseReq.Mode)
	copy(streamSlice[9:], baseReq.Data)
	n, err := c.conn.Write(streamSlice)
	if err != nil {
		return 0, fmt.Errorf("Client Call %v  %v", n, err)
	}

	var slice4byte = make([]byte, 4)
	n, err = io.ReadFull(c.conn, slice4byte)
	if err != nil {
		return 0, fmt.Errorf("Client read packet len err %v  %v", n, err)
	}
	var l = binary.BigEndian.Uint32(slice4byte)

	var t = make([]byte, 1)
	n, err = io.ReadFull(c.conn, t)
	if err != nil {
		return 0, fmt.Errorf("Client read packet type err %v  %v", n, err)
	}

	n, err = io.ReadFull(c.conn, slice4byte)
	if err != nil {
		return 0, fmt.Errorf("Client read packet code err %v  %v", n, err)
	}
	var code = binary.BigEndian.Uint32(slice4byte)

	var payload = make([]byte, l)
	n, err = io.ReadFull(c.conn, payload)
	if err != nil {
		return 0, fmt.Errorf("Client read packet payload err %v  %v", n, err)
	}

	msg.Unmarshal(payload, rsp)
	return code, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
