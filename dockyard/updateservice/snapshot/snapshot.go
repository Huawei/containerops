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

package snapshot

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Callback is a function that a snapshot plugin use after finish the `Process`
// configuration.
type Callback func(id string, output SnapshotOutputInfo) error

type SnapshotInputInfo struct {
	// Snapshot Plugin name.
	// There are two types of Snapshot:
	//   one is simple snapshot, just like 'appv1',calling 'Process' directly.
	//   one is group snapshot, just like 'bycontainer/scanimage', using `scanimage` to 'Process'.
	Name string
	// Dockyard server host address,
	Host string
	// CallbackID, encrytped
	// TODO with timestamp?
	CallbackID string
	// CallbackFunc: if callback is nil, the caller could handle it by itself
	// or the caller must implement calling this in `Process`
	CallbackFunc Callback
	DataProto    string
	// dir/file url of the data
	DataURL string
}

func (info *SnapshotInputInfo) GetName() (string, string) {
	n := strings.SplitN(info.Name, "/", 2)
	if len(n) == 2 {
		return n[0], n[1]
	}
	return n[0], ""
}

// TODO: Better structed
type SnapshotOutputInfo struct {
	Data  []byte
	Error error
}

// UpdateServiceSnapshot represents the snapshot interface
type UpdateServiceSnapshot interface {
	New(info SnapshotInputInfo) (UpdateServiceSnapshot, error)
	// `proto`: `appv1/dockerv1` for example
	Supported(proto string) bool
	Description() string
	Process() error
}

var (
	usSnapshotsLock sync.Mutex
	usSnapshots     = make(map[string]UpdateServiceSnapshot)
)

// RegisterSnapshot provides a way to dynamically register an implementation of a
// snapshot type.
func RegisterSnapshot(name string, f UpdateServiceSnapshot) error {
	if name == "" {
		return errors.New("Could not register a Snapshot with an empty name")
	}
	if f == nil {
		return errors.New("Could not register a nil Snapshot")
	}

	usSnapshotsLock.Lock()
	defer usSnapshotsLock.Unlock()

	if _, alreadyExists := usSnapshots[name]; alreadyExists {
		return fmt.Errorf("Snapshot type '%s' is already registered", name)
	}
	usSnapshots[name] = f
	return nil
}

func UnregisterAllSnapshot() {
	usSnapshotsLock.Lock()
	defer usSnapshotsLock.Unlock()

	for n, _ := range usSnapshots {
		delete(usSnapshots, n)
	}
}

func IsSnapshotSupported(info SnapshotInputInfo) (bool, error) {
	name, _ := info.GetName()
	f, ok := usSnapshots[name]
	if !ok {
		return false, fmt.Errorf("Cannot find plugin :%s", name)
	}

	ok = f.Supported(info.DataProto)
	if !ok {
		return false, fmt.Errorf("Proto %s is not supported by plugin %s", info.DataProto, name)
	}

	return true, nil
}

func ListSnapshotByProto(proto string) (snapshots []string) {
	for n, f := range usSnapshots {
		if f.Supported(proto) {
			snapshots = append(snapshots, n)
		}
	}

	return
}

// NewUpdateServiceSnapshot creates a snapshot interface by an info and a url
func NewUpdateServiceSnapshot(info SnapshotInputInfo) (UpdateServiceSnapshot, error) {
	name, _ := info.GetName()
	f, ok := usSnapshots[name]
	if !ok {
		return nil, fmt.Errorf("Snapshot '%s' not found", name)
	}

	return f.New(info)
}
