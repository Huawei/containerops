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

package km

import (
	"errors"
	"fmt"
	"sync"

	"github.com/containerops/configure"
)

// KeyManager should be seperate from dockyard
// Now only assume that keys are existed in the backend key manager.
// It is up to each implementation to decide whether provides a way
//  to generate key pair automatically.
type KeyManager interface {
	// `url` is the database address or local directory (for example: /tmp/cache)
	New(url string) (KeyManager, error)
	Supported(url string) bool
	// proto: 'app/v1' for example
	// namespace : namespace
	GetPublicKey(proto string, namespace string) ([]byte, error)
	// proto: 'app/v1' for example
	// namespace : namespace
	Sign(proto string, namespace string, data []byte) ([]byte, error)
	// proto: 'app/v1' for example
	// namespace : namespace
	Decrypt(proto string, namespace string, data []byte) ([]byte, error)
	// WARNING! it is dangrous to privide this, so mask it now.
	// In replace, we provides Sign and Decrypt as a service.
	// GetPrivateKey(proto string, namespace string) ([]byte, error)
}

var (
	kmsLock sync.Mutex
	kms     = make(map[string]KeyManager)

	// ErrorsKMNotSupported occurs when the km type is not supported
	ErrorsKMNotSupported = errors.New("key manager type is not supported")
)

// RegisterKeyManager provides a way to dynamically register an implementation of a
// key manager type.
//
// If RegisterKeyManager is called twice with the same name if 'key manager type' is nil,
// or if the name is blank, it panics.
func RegisterKeyManager(name string, f KeyManager) error {
	if name == "" {
		return errors.New("Could not register a KeyManager with an empty name")
	}
	if f == nil {
		return errors.New("Could not register a nil KeyManager")
	}

	kmsLock.Lock()
	defer kmsLock.Unlock()

	if _, alreadyExists := kms[name]; alreadyExists {
		return fmt.Errorf("KeyManager type '%s' is already registered", name)
	}
	kms[name] = f

	return nil
}

// NewKeyManager create a key manager by a url
func NewKeyManager(url string) (KeyManager, error) {
	if url == "" {
		url = configure.GetString("updateserver.keymanager")
	}
	for _, f := range kms {
		if f.Supported(url) {
			return f.New(url)
		}
	}

	return nil, ErrorsKMNotSupported
}
