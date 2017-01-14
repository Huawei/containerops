devops.controller('CreateComponentController', ['$scope','$location', 'componentService',
  function($scope,$location,componentService) {   

    $scope.showEditShell = function(event){
      selectTab(event);
    }

    $scope.showBuildImage = function(event){
      selectTab(event);
    }

    $scope.showRuntimeConfig = function(event){
      selectTab(event);
    }

    $scope.showDebug = function(event){
      selectTab(event);
    }

    function selectTab(event){
      $(".component-tab").removeClass("component-tab-active");
      $(event.currentTarget).addClass("component-tab-active");
    }

    $scope.baseOrAdvanced = function(){
      if($scope.component.use_advanced){
        $scope.component.service = {};
        $scope.component.pod = {};
      }else{
        $scope.component.service = $.extend(true,{},componentService.base_service);
        $scope.component.pod = $.extend(true,{},componentService.base_pod);
      }
    }

    function init(){
      $scope.component = $.extend(true,{},componentService.component);
      $scope.runtimeTab = 1;

      if(!$scope.component.use_advanced){
        $scope.component.service = $.extend(true,{},componentService.base_service);
        $scope.component.pod = $.extend(true,{},componentService.base_pod);
      } 
    }

    init();
}]);