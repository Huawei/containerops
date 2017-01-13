devops.controller('ComponentController', ['$scope','$location','componentService', 'notifyService', 'loading', 'more',
  function($scope,$location,componentService,notifyService,loading,more) {   

    $scope.filter = {
        "name" : "",
        "filter" : ""
    }

    $scope.getComponents = function(){
      var promise = componentService.getComponents("","",true,10,3,0);
      promise.done(function(data){
          loading.hide();
          $scope.components = data.list;
      });
      promise.fail(function(xhr,status,error){
          loading.hide();
          if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.common) {
              notifyService.notify(xhr.responseJSON.common.error_code + " : " + xhr.responseJSON.common.message,"error");
          }else if(xhr.statusText != "abort"){
              notifyService.notify("Server is unreachable","error");
          }
      }); 
    }

    $scope.getComponents();

    more.show(function(){
      alert("bottom!!");
    });

    $scope.components = [
        {
          "name" : "component1",
          "versions" : [
            {
              "id" : 1,
              "version" : "v1.0"
            },
            {
              "id" : 2,
              "version" : "v2.0"
            },
            {
              "id" : 3,
              "version" : "v3.0"
            }
          ]
        },
        {
          "name" : "component2",
          "versions" : [
            {
              "id" : 4,
              "version" : "v1.0"
            },
            {
              "id" : 5,
              "version" : "v2.0"
            }
          ]
        },
        {
          "name" : "component3",
          "versions" : [
            {
              "id" : 6,
              "version" : "1.0"
            },
            {
              "id" : 7,
              "version" : "2.0"
            },
            {
              "id" : 8,
              "version" : "2.0"
            }
          ]
        }
    ]
}]);