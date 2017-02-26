define(["app","services/diagram/main"], function(app) {
    app.controllerProvider.register('WorkflowElementController', ['$scope', '$state', 'notifyService', 'diagramService', function($scope, $state, notifyService, diagramService) {
        $scope.workflows = [{"name":"workflow1","versions":["v1.0","v1.1"]}]
        // $scope.showCreateWorkflow = function(){
        // 	$state.go("workflow.detail");
        // };
        // $scope.showMoreVersion = function(verionname){
        // 	notifyService.notify("no more workflow","info")
        // };

        $scope.workflowData = diagramService.workflowData;
        $scope.currentStageInfo = $scope.workflowData[diagramService.currentStageIndex];

    }]);
})
