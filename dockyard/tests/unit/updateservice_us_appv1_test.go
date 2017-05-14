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
package unittest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/us"
)

// TestUSABasic tests the basic functions
func TestUSABasic(t *testing.T) {
	var appV1 us.UpdateServiceAppV1

	ok := appV1.Supported("appV1")
	assert.Equal(t, ok, true, "Fail to get supported status")
	ok = appV1.Supported("appInvalid")
	assert.Equal(t, ok, false, "Fail to get supported status")
}
