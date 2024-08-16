/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

const (
	// GenericService name
	GenericService = "$GenericService" // private as "$"
	// GenericMethod name
	GenericMethod = "$GenericCall"
)

// ServiceInfo to record meta info of service
type ServiceInfo struct {
	// The name of the service. For generic services, it is always the constant `GenericService`.
	ServiceName string

	// HandlerType is the type value of a request serviceImpl from the generated code.
	HandlerType interface{}

	// Methods contains the meta information of methods supported by the service.
	// For generic service, there is only one method named by the constant `GenericMethod`.
	Methods map[string]MethodInfo
	// GenericMethod returns a MethodInfo for the given name.
	// It is used by generic calls only.
	GenericMethod func(name string) MethodInfo
}

// MethodInfo gets MethodInfo.
func (i *ServiceInfo) MethodInfo(name string) MethodInfo {
	if i.ServiceName == GenericService {
		if i.GenericMethod != nil {
			return i.GenericMethod(name)
		}
		return i.Methods[GenericMethod]
	}
	return i.Methods[name]
}

func NewServiceInfo(serviceName string, HandlerType interface{}, Methods map[string]MethodInfo) *ServiceInfo {
	return &ServiceInfo{
		ServiceName: serviceName,
		HandlerType: HandlerType,
		Methods:     Methods,
	}
}
