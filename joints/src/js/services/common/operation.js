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
define(['services/module'], function(commonServiceModule) {
    commonServiceModule.factory("loading", function() {
            function show() {
                $(".loading").removeClass("hide");
            }

            function hide() {
                $(".loading").addClass("hide");
            }

            return {
                "show": show,
                "hide": hide
            }
        })
        .factory("more", ['$location', function($location) {
            return {
                show: function(next) {
                    $(window).scroll(function() {
                        if ($location.path() == "/component") {
                            if ($(window).scrollTop() + $(window).height() == $(document).height()) {
                                next();
                            }
                        }
                    });
                }
            }
        }])
        .factory("notifyService", function() {
            Messenger.options = { extraClasses: 'messenger-fixed messenger-theme-future messenger-on-top'};
            return {
                notify: function(msg, type, showtime) {
                    Messenger().post({
                        "message": msg,
                        "type": type,
                        /* success, error, info*/
                        "showCloseButton": true,
                        "hideAfter": showtime ? showtime : 3
                    });
                }
            }
        })
})
