/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package main

import "testing"

func Test_parse_env(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name       string
		args       args
		wantURI    string
		wantAction string
		wantErr    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURI, gotAction, err := parse_env(tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse_env() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotURI != tt.wantURI {
				t.Errorf("parse_env() gotURI = %v, want %v", gotURI, tt.wantURI)
			}
			if gotAction != tt.wantAction {
				t.Errorf("parse_env() gotAction = %v, want %v", gotAction, tt.wantAction)
			}
		})
	}
}
