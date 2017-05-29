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

define(['services/module'], function(commonServiceModule) {
    commonServiceModule.factory("apiService", ["notifyService", "loading", function(notifyService, loading) {
        var pendingPromise;
        // abort
        function beforeApiInvocation(skipAbort) {
            if (!skipAbort) {
                _.each(pendingPromise, function(promise) {
                    promise.abort();
                });
                pendingPromise = [];
            }
            loading.show();
        }

        function addPromise(promise) {
            pendingPromise.push(promise);
        }

        function failToCall(error) {
            loading.hide();
            if (!_.isUndefined(error) && error.common) {
                notifyService.notify(error.common.error_code + " : " + error.common.message, "error");
            } else {
                notifyService.notify("Server is unreachable", "error");
            }
        }

        function successToCall(msg) {
            loading.hide();
            if (!_.isUndefined(msg) && msg.common) {
                notifyService.notify(msg.common.message, "success");
            }
        }

        return {
            "beforeApiInvocation": beforeApiInvocation,
            "addPromise": addPromise,
            "failToCall": failToCall,
            "successToCall": successToCall
        }
    }])
})
