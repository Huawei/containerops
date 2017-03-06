define(["app","services/diagram/main"], function(app) {
    app.controllerProvider.register('WorkflowElementController', ['$scope', '$rootScope', '$state', 'notifyService', 'diagramService', function($scope, $rootScope, $state, notifyService, diagramService) {
        $scope.workflows = [{"name":"workflow1","versions":["v1.0","v1.1"]}]
        // $scope.showCreateWorkflow = function(){
        // 	$state.go("workflow.detail");
        // };
        // $scope.showMoreVersion = function(verionname){
        // 	notifyService.notify("no more workflow","info")
        // };
        var currentStageIndex = diagramService.currentStageIndex;
        $scope.workflowData = diagramService.workflowData;
        $scope.currentStageInfo = $scope.workflowData[currentStageIndex];
        $scope.isShowPage = true;

        $scope.redraw = function(){
            diagramService.resetWorkflowData($scope.workflowData);
            $rootScope.drawWorkflow();
        };

        $scope.stageEvent = {
            delete: function(){
                $scope.workflowData.splice(currentStageIndex,1);
                $scope.redraw();
                $scope.isShowPage = false;
            }
        };

    }]);
})
