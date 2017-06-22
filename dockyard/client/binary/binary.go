/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package binary

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Upload binary file to the Dockyard service.
func UploadBinaryFile(filePath, domain, namespace, repository, tag string) error {
	if f, err := os.Open(filePath); err != nil {
		return err
	} else {
		defer f.Close()

		if req, err := http.NewRequest(http.MethodPut,
			fmt.Sprintf("https://%s/binary/v1/%s/%s/binary/%s/%s",
				domain, namespace, repository, filepath.Base(filePath), tag), f); err != nil {
			return err
		} else {
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			if resp, err := client.Do(req); err != nil {
				return err
			} else {
				defer resp.Body.Close()

				switch resp.StatusCode {
				case http.StatusOK:
					return nil
				case http.StatusBadRequest:
					return fmt.Errorf("Binary upload failed.")
				case http.StatusUnauthorized:
					return fmt.Errorf("Action unauthorized.")
				default:
					return fmt.Errorf("Unknown error.")
				}
			}
		}
	}

	return nil
}

// Download binary file to the local.
func DownloadBinaryFile(domain, namespace, repository, filename, tag, filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	if f, err := os.Create(filePath); err != nil {
		return err
	} else {
		defer f.Close()

		if resp, err := http.Get(fmt.Sprintf("https://%s/binary/v1/%s/%s/binary/%s/%s",
			domain, namespace, repository, filename, tag)); err != nil {
			return err
		} else {
			defer resp.Body.Close()

			if _, err := io.Copy(f, resp.Body); err != nil {
				return err
			}
		}
	}

	return nil
}
