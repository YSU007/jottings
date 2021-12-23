package path

import "fmt"

// IServiceDiscoverPath 服务发现路径
type IServiceDiscoverPath interface {
	GetPath() string
	GetCompany() string
	GetVersion() string
	GetServiceName() string
}

//DefaultDiscoverPath 默认path struct
type DefaultDiscoverPath struct {
	Company     string
	Version     string
	ServiceName string
}

func (d DefaultDiscoverPath) GetPath() string {
	return fmt.Sprintf("/%s/%s/%s/", d.GetCompany(), d.GetVersion(), d.GetServiceName())
}

func (d DefaultDiscoverPath) GetCompany() string {
	return d.Company
}

func (d DefaultDiscoverPath) GetVersion() string {
	return d.Version
}

func (d DefaultDiscoverPath) GetServiceName() string {
	return d.ServiceName
}
