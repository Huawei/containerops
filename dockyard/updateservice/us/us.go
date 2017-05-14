/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package us

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Huawei/dockyard/utils"
)

// UpdateService represents the update service interface
type UpdateService interface {
	Supported(protocal string) bool
	New(protocal string, storageURL string, kmURL string) (UpdateService, error)
	List(nr string) ([]string, error)
	GetPublicKey(nr string) ([]byte, error)
	GetMeta(nr string) ([]byte, error)
	GetMetaSign(nr string) ([]byte, error)
	Get(nr string, name string) ([]byte, error)
	// Return the path of local storage or key of object store
	Put(nr string, name string, data []byte, method utils.EncryptMethod) (string, error)
	Delete(nr string, name string) error
}

var (
	ussLock sync.Mutex
	uss     = make(map[string]UpdateService)

	// ErrorsUSPNotSupported occurs when a protocal is not supported
	ErrorsUSPNotSupported = errors.New("protocal is not supported")
)

// Register provides a way to dynamically register an implementation of a
// protocal.
//
// If Register is called twice with the same name if 'protocal' is nil,
// or if the name is blank, it panics.
func Register(name string, f UpdateService) error {
	if name == "" {
		return errors.New("Could not register a  with an empty name")
	}
	if f == nil {
		return errors.New("Could not register a nil ")
	}

	ussLock.Lock()
	defer ussLock.Unlock()

	if _, alreadyExists := uss[name]; alreadyExists {
		return errors.New(fmt.Sprintf(" type '%s' is already registered", name))
	}
	uss[name] = f

	return nil
}

// NewUpdateService create a update service protocal interface by a protocal type
func NewUpdateService(protocal string, storageURL string, kmURL string) (UpdateService, error) {
	for _, f := range uss {
		if f.Supported(protocal) {
			return f.New(protocal, storageURL, kmURL)
		}
	}

	return nil, ErrorsUSPNotSupported
}
