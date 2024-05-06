package main

import "fmt"

// IService 定义了服务接口
type IService interface {
	Serve()
}

// RealService 是实现 IService 接口的真实服务
type RealService struct{}

// Serve 是 RealService 的方法，表示提供服务
func (s *RealService) Serve() {
	fmt.Println("RealService: 提供服务")
}

// ProxyService 是代理服务，它包含了对真实服务的引用
type ProxyService struct {
	realService *RealService
}

// Serve 是 ProxyService 的方法，它在调用真实服务之前可以执行一些操作
func (p *ProxyService) Serve() {
	fmt.Println("ProxyService: 在调用真实服务之前的操作")
	p.realService.Serve()
}

func main() {
	realService := &RealService{}
	proxyService := &ProxyService{realService: realService}

	// 通过代理服务来使用真实服务
	proxyService.Serve()
}
