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

define(['app'], function(app) {
    app.provide.factory("componentCheck", ["notifyService", function(notifyService) {

        // runtime config
        function name() {
            if ($('#component-name-form').length > 0) {
                return $('#component-name-form').parsley().validate();
            } else {
                return true;
            }
        }

        function version() {
            if ($('#component-version-form').length > 0) {
                return $('#component-version-form').parsley().validate();
            } else {
                return true;
            }
        }

        function existingimage() {
            return $('#component-form').parsley().validate();
        }

        function basesetting() {
            return $('#base-setting-form').parsley().validate();
        }

        function advancedsetting() {
            return $('#advanced-setting-form').parsley().validate();
        }

        function env() {
            return $('#component-envs').parsley().validate();
        }

        var imageSetting = {
            "from": function() {
                return $('#base-image-form').parsley().validate();
            },
            "events": {
                "componentStart": function(componentData) {
                    return _.isEmpty(componentData.imageSetting.events.componentStart);
                },
                "componentResult": function(componentData) {
                    return _.isEmpty(componentData.imageSetting.events.componentResult);
                },
                "componentStop": function(componentData) {
                    return _.isEmpty(componentData.imageSetting.events.componentStop);
                }
            },
            "push": function() {
                return $('#push-image-form').parsley().validate();
            }
        }

        function go(componentData) {
            var check_result = true;

            if (!name()) {
                check_result = false;
                notifyService.notify("Component name is required.", "error");
                return check_result;
            }

            if (!version()) {
                check_result = false;
                notifyService.notify("Component version is required.", "error");
                return check_result;
            }

            if (!existingimage()) {
                check_result = false;
                notifyService.notify("Repository name of existing image is required.", "error");
                return check_result;
            }

            if (imageSetting.events.componentStart(componentData)) {
                check_result = false;
                notifyService.notify("Script of component start event is required.", "error");
                return check_result;
            }
            if (imageSetting.events.componentResult(componentData)) {
                check_result = false;
                notifyService.notify("Script of component result event is required.", "error");
                return check_result;
            }
            if (imageSetting.events.componentStop(componentData)) {
                check_result = false;
                notifyService.notify("Script of component stop event is required.", "error");
                return check_result;
            }
            if (!imageSetting.from()) {
                check_result = false;
                notifyService.notify("Base image is not complete.", "error");
                return check_result;
            }
            if (!imageSetting.push()) {
                check_result = false;
                notifyService.notify("Push image is not complete.", "error");
                return check_result;
            }

            if (!componentData.useAdvanced) {
                if (!basesetting()) {
                    check_result = false;
                    notifyService.notify("Kubernetes base setting is not complete.", "error");
                    return check_result;
                }
            } else {
                if (!advancedsetting()) {
                    check_result = false;
                    notifyService.notify("Kubernetes advanced setting is not complete.", "error");
                    return check_result;
                }
            }

            if (componentData.env.length > 0) {
                if (!env()) {
                    check_result = false;
                    notifyService.notify("Component env is not complete.", "error");
                    return check_result;
                }
            }

            return check_result;
        }

        var tabcheck = {
            "runtime": function(componentData) {
                var result = true;

                if (!existingimage()) {
                    result = false;
                }

                if (!componentData.useAdvanced) {
                    if (!basesetting()) {
                        result = false;
                    }
                } else {
                    if (!advancedsetting()) {
                        result = false;
                    }
                }

                if (componentData.env.length > 0) {
                    if (!env()) {
                        result = false;
                    }
                }

                return result;
            },
            "editshell": function(componentData) {
                var result = true;
                if (imageSetting.events.componentStart(componentData)) {
                    result = false;
                } else if (imageSetting.events.componentResult(componentData)) {
                    result = false;
                } else if (imageSetting.events.componentStop(componentData)) {
                    result = false;
                }
                return result;
            },
            "buildimage": function() {
                var result = imageSetting.from() && imageSetting.push();
                return result;
            }
        }

        return {
            "version": version,
            "go": go,
            "tabcheck": tabcheck
        }
    }])
})
