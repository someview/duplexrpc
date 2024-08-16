package service

import (
	"errors"
)

type Service struct {
	svcInfo     *ServiceInfo
	serviceImpl interface{}
}

func NewService(svcInfo *ServiceInfo, handler interface{}) Service {
	return Service{svcInfo: svcInfo, serviceImpl: handler}
}

func (s *Service) GetMethodInfoAndSvcImpl(methodName string) (MethodInfo, any) {
	return s.svcInfo.MethodInfo(methodName), s.serviceImpl
}

type Manager map[string]Service

func (s Manager) AddService(svc Service) error {
	serviceName := svc.svcInfo.ServiceName
	if _, ok := s[serviceName]; ok {
		return errors.New("service name has been registered ")
	}
	s[serviceName] = svc
	return nil
}

func (s Manager) GetSvcInfoMap() map[string]*ServiceInfo {
	svcInfoMap := map[string]*ServiceInfo{}
	for name, svc := range s {
		svcInfoMap[name] = svc.svcInfo
	}
	return svcInfoMap
}

// return service impl  methodInfo

func (s Manager) GetService(svcName string) (svc Service, ok bool) {
	svc, ok = s[svcName]
	return
}
