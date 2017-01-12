devops.controller('ComponentController', ['$scope','$location','componentService', 'notifyService',
  function($scope,$location,componentService,notifyService) {  

  $scope.getComponents = function(){
    var promise = componentService.getComponents();
    promise.done(function(data){
        // loading.hide();
        $scope.components = data.list;
    });
    promise.fail(function(xhr,status,error){
        // loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notifyService.notify(xhr.responseJSON.errMsg,"error");
        }else if(xhr.statusText != "abort"){
            notifyService.notify("Server is unreachable","error");
        }
    }); 
  }

  $scope.getComponents();
}]);