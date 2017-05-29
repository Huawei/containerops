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

define(["app","services/component/main"], function(app) {
    app.controllerProvider.register('ComponentController', ['$scope', '$state', 'componentService', 'notifyService', 'loading',
        'apiService', 'utilService',
        function($scope, $state, componentService, notifyService, loading, apiService, utilService) {

            function getOffset(type, name) {
                if (type == "component") {
                    return $scope.components.length;
                } else {
                    return _.find($scope.components, function(component) {
                        return component.name == name;
                    }).versions.length;
                }
            }

            $scope.showMoreComponent = function() {
                var promise = componentService.getComponents($scope.filter.name, $scope.filter.version, true, $scope.pageNum, $scope.versionNum, getOffset("component"));
                promise.done(function(data) {
                    loading.hide();
                    appendComponents(data.components);
                });
                promise.fail(function(xhr, status, error) {
                    apiService.failToCall(xhr.responseJSON);
                });
            }

            function appendComponents(data) {
                var components = utilService.componentDataTransfer(data);
                $scope.components = $scope.components.concat(components);
                $scope.$apply();
            }

            $scope.showMoreVersion = function(componentName) {
                var promise = componentService.getComponents(componentName, "", false, $scope.pageNum, $scope.versionNum, getOffset("version", componentName));
                promise.done(function(data) {
                    loading.hide();
                    appendVersions(data.components, componentName);
                });
                promise.fail(function(xhr, status, error) {
                    apiService.failToCall(xhr.responseJSON);
                });
            }

            function appendVersions(data, componentName) {
                var target = _.find($scope.components, function(component) {
                    return component.name == componentName;
                });
                _.each(data, function(item) {
                    var version = {
                        "id": item.id,
                        "version": item.version
                    }
                    target.versions.push(version);
                })
                $scope.$apply();
            }

            $scope.showNewComponent = function() {
                $state.go("component.create");
            }

            $scope.getComponents = function(type) {
                var promise = componentService.getComponents($scope.filter.name, $scope.filter.version, true, $scope.pageNum, $scope.versionNum, 0);
                promise.done(function(data) {
                    loading.hide();
                    $scope.components = utilService.componentDataTransfer(data.components);
                    $scope.dataReady = true;
                    if (type == "init" && $scope.components.length > 0) {
                        $scope.nodata = false;
                    }
                    $scope.$apply();
                });
                promise.fail(function(xhr, status, error) {
                    $scope.dataReady = true;
                    $scope.components = [];
                    $scope.$apply();
                    apiService.failToCall(xhr.responseJSON);
                });
            }

            $scope.showComponentDetail = function(id) {
                $state.go("component.detail",{"id" : id});
            }

            $scope.confirmDeleteComponent = function(id,componentName,version){
                var actions = [
                    {
                        "name" : "delete",
                        "label" : "Yes",
                        "action" : function action(){
                            deleteComponent(id);
                        }
                    },
                    {
                        "name" : "cancel",
                        "label" : "No",
                        "action" : function action(){
                        }
                    }
                ];
                notifyService.confirm("Are you sure to delete the component " + componentName + " " + version + "?", "info", actions);
            }

            function deleteComponent(id){
                var promise = componentService.deleteComponent(id);
                promise.done(function(data) {
                    loading.hide();
                    $scope.getComponents("init");
                });
                promise.fail(function(xhr, status, error) {
                    apiService.failToCall(xhr.responseJSON);
                });
            }

            function init() {
                $scope.filter = {
                    "name": "",
                    "version": ""
                }

                $scope.pageNum = 10;
                $scope.versionNum = 3;

                $scope.components = [];
                $scope.nodata = true;

                $scope.dataReady = false;

                $scope.getComponents("init");
            }

            init();
        }
    ]);
})
