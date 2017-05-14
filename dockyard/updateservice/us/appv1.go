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
	"github.com/Huawei/dockyard/updateservice/storage"
	"github.com/Huawei/dockyard/utils"
)

const (
	appV1Prefix = "appV1"
	appV1Proto  = "app/v1"
)

// UpdateServiceAppV1 is the appV1 implementation of the update service proto
type UpdateServiceAppV1 struct {
	s storage.UpdateServiceStorage
}

func init() {
	Register(appV1Prefix, &UpdateServiceAppV1{})
}

// Supported checks if a proto is 'appV1'
func (app *UpdateServiceAppV1) Supported(proto string) bool {
	return proto == appV1Prefix
}

// New creates a update service interface by an appV1 proto
func (app *UpdateServiceAppV1) New(proto string, storageURL string, kmURL string) (UpdateService, error) {
	if proto != appV1Prefix {
		return nil, ErrorsUSPNotSupported
	}

	var err error
	app.s, err = storage.NewUpdateServiceStorage(storageURL, kmURL)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Put adds a appV1 file to a repository
func (app *UpdateServiceAppV1) Put(nr, name string, data []byte, method utils.EncryptMethod) (string, error) {
	key := nr + "/" + name
	return app.s.Put(appV1Proto, key, data, method)
}

// Delete removes a appV1 file from a repository
func (app *UpdateServiceAppV1) Delete(nr, name string) error {
	key := nr + "/" + name
	return app.s.Delete(appV1Proto, key)
}

// Get gets the appV1 file data of a repository
func (app *UpdateServiceAppV1) Get(nr, name string) ([]byte, error) {
	key := nr + "/" + name
	return app.s.Get(appV1Proto, key)
}

// List lists the applications of a repository
func (app *UpdateServiceAppV1) List(nr string) ([]string, error) {
	return app.s.List(appV1Proto, nr)
}

// GetPublicKey returns the public key data of a repository
func (app *UpdateServiceAppV1) GetPublicKey(namespace string) ([]byte, error) {
	return app.s.GetPublicKey(appV1Proto, namespace)
}

// GetMeta returns the meta data of a repository
func (app *UpdateServiceAppV1) GetMeta(nr string) ([]byte, error) {
	return app.s.GetMeta(appV1Proto, nr)
}

// GetMetaSign returns the meta signature data of a repository
func (app *UpdateServiceAppV1) GetMetaSign(nr string) ([]byte, error) {
	return app.s.GetMetaSign(appV1Proto, nr)
}
