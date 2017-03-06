define(["app","services/diagram/main"], function(app) {
    app.controllerProvider.register('ActionDetailController', ['$scope', '$rootScope', '$state', '$stateParams', 'notifyService', 'diagramService', function($scope, $rootScope, $state, $stateParams, notifyService, diagramService) {
        $scope.workflowData = diagramService.workflowData;
        $scope.isShowPage = true;

        var currentStageIndex = diagramService.currentStageIndex;
        var currentActionIndex = diagramService.currentActionIndex;
        var id = $stateParams.id;


        $scope.redraw = function(){
            diagramService.resetWorkflowData($scope.workflowData);
            $rootScope.drawWorkflow();
        };

        $scope.getCurrentAction = function(id){
            angular.forEach($scope.workflowData,function(d,i){
                angular.forEach(d.actions,function(a,ai){
                    if(a.id === id){
                        currentStageIndex = i;
                        currentActionIndex = ai;
                        a.isChosed = true;
                        $scope.redraw();
                    }
                });
            });
        };

        if(id){
            $scope.getCurrentAction(id);
        };

        // console.log(currentStageIndex)

        $scope.currentActionInfo = $scope.workflowData[currentStageIndex]['actions'][currentActionIndex];

        $scope.actionEvent = {
            delete: function(){
                $scope.workflowData[currentStageIndex]['actions'].splice(currentActionIndex,1);
                $scope.redraw();
                $scope.isShowPage = false;
            }
        };

        $scope.componentEvent = {
            delete: function(index){
                $scope.currentActionInfo.components.splice(index,1);
                $scope.redraw();
            },
            add: function(){
                $state.go("workflow.create.addComponent");
            },
            edit: function(index,item){
                diagramService.currentComponentIndex = index;
                $state.go('workflow.create.editComponent',{id:item.id});
            }
        }
    }]);
})
