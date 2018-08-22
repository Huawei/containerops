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

package template

import (
	"regexp"
)

//getTemplateContent get the template content by given version, if no content
//for the version, try its parenet version(like, 1.7.8 -> 1.7)
func getTemplateContent(templates map[string]string, version string) string {
	tplContent := templates[version]
	if tplContent == "" {
		parentVersion := parentVersionPattern.FindString(version)
		// TODO Log the Warnning
		// objects.WriteLog(fmt.Sprintf("No ca template for version %s, trying to use parent version %s", version, parentVersion), stdout, timestamp, d)

		tplContent = templates[parentVersion]
		if tplContent == "" {
			// return files, fmt.Errorf("No template for version %s or %s", version, parentVersion)
			// TODO Log the error
		}
	}
	return tplContent
}

var parentVersionPattern *regexp.Regexp = regexp.MustCompile(`[a-z]+-\d+\.\d+`)

func getParentVersion(version string) string {
	parentVersion := parentVersionPattern.FindString(version)
	return parentVersion
}
