package main

func main() {
}

// AccountInterface ----------------------------------------------------------------------------------------------------
type AccountInterface interface {
	GetAccountId() string
}

// RequestInterface ResponseInterface ----------------------------------------------------------------------------------------------------
type RequestInterface interface {
	GetMode() uint32
}

type ResponseInterface interface {
	GetCode() uint32
}

// RouterInterface ----------------------------------------------------------------------------------------------------
type RouterInterface interface {
	RegHandle(mode uint32, handleInterface HandleInterface)
	HandleServe(a *AccountInterface, req *RequestInterface, rsp *ResponseInterface)
}

type HandleInterface interface {
	Serve(a *AccountInterface, req *RequestInterface, rsp *ResponseInterface)
}

type MappingRouter map[uint32]HandleInterface

func (r *MappingRouter) RegHandle(mode uint32, handleInterface HandleInterface) {
	(*r)[mode] = handleInterface
}

func (r *MappingRouter) HandleServe(a *AccountInterface, req *RequestInterface, rsp *ResponseInterface) {
	var mode = (*req).GetMode()
	(*r)[mode].Serve(a, req, rsp)
}

var DefRouter = make(MappingRouter)

func RegHandles() {
}

// AccountImpl ----------------------------------------------------------------------------------------------------
type AccountImpl struct {
	AccountId string
}

func (a *AccountImpl) GetAccountId() string {
	return a.AccountId
}

// RequestBase ResponseBase ----------------------------------------------------------------------------------------------------
type RequestBase struct {
	Mode uint32
}

func (r *RequestBase) GetMode() uint32 {
	return r.Mode
}

type ResponseBase struct {
	Code uint32
}

func (r *ResponseBase) GetCode() uint32 {
	return r.Code
}

// TestHandle ----------------------------------------------------------------------------------------------------
func (a *AccountImpl) TestHandle(req *RequestBase) (rsp *ResponseBase) {
	return nil
}
