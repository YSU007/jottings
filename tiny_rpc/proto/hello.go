package proto

// HelloReq HelloRsp ----------------------------------------------------------------------------------------------------
type HelloReq struct {
	HelloMsg string `json:"hello_msg,omitempty"`
}

type HelloRsp struct {
	ReplyMsg string `json:"reply_msg,omitempty"`
}
