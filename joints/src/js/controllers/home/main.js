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

define(["app"], function(app) {
    app.controllerProvider.register('HomeController', ['$scope', '$location', function($scope, $location) {

        $scope.navigators = [{
            "name": "Workflow",
            "href": "/workflow",
            "ngclass": "active",
            "icon": "fa fa-desktop"
        }, {
            "name": "Component",
            "href": "/component",
            "ngclass": "",
            "icon": "fa fa-cube"
        }, {
            "name": "History",
            "href": "/history",
            "ngclass": "",
            "icon": "fa fa-history"
        }, {
            "name": "System Setting",
            "href": "/setting",
            "ngclass": "",
            "icon": "fa fa-cog"
        }];

        $scope.chooseNav = function(name) {
            _.each($scope.navigators, function(item) {
                if (item.name != name) {
                    item.ngclass = "";
                } else {
                    item.ngclass = "active";
                }
            })
        }

        function initializeNav() {
            var path = $location.path();
            _.each($scope.navigators, function(item) {
                if (path != "") {
                    if (path.indexOf(item.href) == 0) {
                        item.ngclass = "active";
                    } else {
                        item.ngclass = "";
                    }
                }
            })
        }

        initializeNav();

    }]);
})
