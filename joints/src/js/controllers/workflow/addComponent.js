define(["app","services/diagram/main","services/component/main"], function(app) {
    app.controllerProvider.register('AddComponentController', ['$scope', '$rootScope', '$state', 'notifyService', 'diagramService', 'componentService', 'utilService', 'apiService', 'loading', function($scope, $rootScope, $state, notifyService, diagramService, componentService, utilService, apiService, loading) {
        $scope.workflowData = diagramService.workflowData;
        var currentStageIndex = diagramService.currentStageIndex;
        var currentActionIndex = diagramService.currentActionIndex;
        $scope.currentActionInfo = $scope.workflowData[currentStageIndex]['actions'][currentActionIndex];
        // $scope.componentEvent = {
        //     delete: function(index){
        //         $scope.currentActionInfo.components.splice(index,1);
        //     },
        //     add: function(){
        //         $state.go("workflow.create.action",{"id": d.id});
        //     },
        //     edit: function(index){
        //         console.log($scope.currentActionInfo.components[index])
        //     }
        // }
         $scope.showMoreVersion = function(componentName) {
            var promise = componentService.getComponents(componentName, "", false, $scope.pageNum, $scope.versionNum, getOffset("version", componentName));
            promise.done(function(data) {
                loading.hide();
                appendVersions(data.components, componentName);
            });
            promise.fail(function(xhr, status, error) {
                apiService.failToCall(xhr.responseJSON);
            });
        };

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
        };

         function getOffset(type, name) {
            if (type == "component") {
                return $scope.components.length;
            } else {
                return _.find($scope.components, function(component) {
                    return component.name == name;
                }).versions.length;
            }
        };

        $scope.getComponents = function(type) {
            var promise = componentService.getComponents($scope.filter.name, $scope.filter.version, true, $scope.pageNum, $scope.versionNum, 0);
            promise.done(function(data) {
                loading.hide();
                $scope.components = utilService.componentDataTransfer(data.components);
                $scope.addDefultTimes($scope.components);
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
        };

        $scope.addDefultTimes = function(data){
            angular.forEach(data,function(d,i){
                angular.forEach(d.versions,function(v,vi){
                    v.importTimes = 0;
                })
            })
        };

        $scope.addImportComponent = function(id){
            angular.forEach($scope.components,function(d,i){
                angular.forEach(d.versions,function(v,vi){
                    if(v.id === id){
                        v.importTimes++;
                        getComponent(id);
                    }
                })
            })
        };

        function getComponent(id) {
            var promise = componentService.getComponent(id);
            promise.done(function(data) {
                loading.hide();
                var component = data.component;
                    component.uuid = uuid.v1();
                $scope.importComponents.push(component);

            });
            promise.fail(function(xhr, status, error) {
                apiService.failToCall(xhr.responseJSON);
            });
        };

        $scope.closeComponents = function(){
            var workflowData = angular.copy(diagramService.workflowData);
            var currentAction = workflowData[diagramService.currentStageIndex]['actions'][diagramService.currentActionIndex];
            currentAction.components = currentAction.components.concat($scope.importComponents);
            diagramService.resetWorkflowData(workflowData);
            $rootScope.drawWorkflow(workflowData);
            $state.go("workflow.create.action",{id:currentAction.id});
        };


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

            $scope.importComponents = [];

            $scope.getComponents("init");
        };
        init();

    }]);
})
