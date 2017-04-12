package v1

import "github.com/lastbackend/lastbackend/pkg/apis/types"

func NewService(obj *types.Service) *Service {
	return New(obj)
}

func NewServiceList(obj *types.ServiceList) *ServiceList {
	return NewList(obj)
}
