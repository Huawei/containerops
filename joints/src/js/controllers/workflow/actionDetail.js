define(["app","services/diagram/main"], function(app) {
    app.controllerProvider.register('ActionDetailController', ['$scope', '$state', 'notifyService', 'diagramService', function($scope, $state, notifyService, diagramService) {
        $scope.workflowData = diagramService.workflowData;
        var currentStageIndex = diagramService.currentStageIndex;
        var currentActionIndex = diagramService.currentActionIndex;
        $scope.currentActionInfo = $scope.workflowData[currentStageIndex]['actions'][currentActionIndex];
        $scope.componentEvent = {
            delete: function(index){
                $scope.currentActionInfo.components.splice(index,1);
            },
            add: function(){
                $state.go("workflow.create.addComponent");
            },
            edit: function(index){
                console.log($scope.currentActionInfo.components[index])
            }
        }
    }]);
})
