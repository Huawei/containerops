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

package utils

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"
)

// Meta represents the meta info of a repository
type Meta struct {
	Items   []MetaItem
	Updated time.Time
}

// MetaItem represents the meta information of a repository app/vm/image
type MetaItem struct {
	Name   string
	Hash   string
	Method EncryptMethod

	Created time.Time
	Expired time.Time
}

const (
	//The default life circle for a software is half a year
	defaultLifecircle = time.Hour * 24 * 180
)

func (a Meta) Before(b Meta) bool {
	return a.Updated.Before(b.Updated)
}

// GenerateMetaItem generates a meta data by a file name and file content
func GenerateMetaItem(file string, contentByte []byte) (meta MetaItem) {
	meta.Name = file
	meta.Hash = fmt.Sprintf("%x", sha1.Sum(contentByte))
	meta.Created = time.Now()
	meta.Expired = meta.Created.Add(defaultLifecircle)
	return
}

// GetHash get the hash string of a file
func (a MetaItem) GetHash() string {
	return a.Hash
}

func (a *MetaItem) SetEncryption(method EncryptMethod) {
	a.Method = method
}

func (a MetaItem) GetEncryption() EncryptMethod {
	return a.Method
}

// IsExpired tells if an application is expired
func (a MetaItem) IsExpired() bool {
	//FIXME: read time from time server?
	return a.Expired.Before(time.Now())
}

// GetCreated returns the created time of an application
func (a MetaItem) GetCreated() time.Time {
	return a.Created
}

// SetCreated set the created time of an application
func (a *MetaItem) SetCreated(t time.Time) {
	a.Created = t
}

// GetExpired get the expired time of an application
func (a MetaItem) GetExpired() time.Time {
	return a.Expired
}

// SetExpired set the expired time of an application
func (a *MetaItem) SetExpired(t time.Time) {
	a.Expired = t
}

// Compare checks if two meta is the same
func (a MetaItem) Compare(b MetaItem) int {
	if reflect.DeepEqual(a, b) {
		return 0
	}

	if a.Created.Before(b.Created) {
		return -1
	}

	return 1
}
