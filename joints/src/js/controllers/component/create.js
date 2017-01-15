devops.controller('CreateComponentController', ['$scope','$location', 'componentService', 'componentIO',
  function($scope,$location,componentService,componentIO) {   

    // tabs control
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

    // runtime tab all functions below
    $scope.changeRuntimeTab = function(index){
      $scope.runtimeTab = index;
      if($scope.runtimeTab == 2){
        setTimeout(function(){
          componentIO.init($scope.component);
        },1000)
      }
    }

    $scope.baseOrAdvanced = function(){
      if($scope.component.use_advanced){
        $scope.component.service = {};
        $scope.component.pod = {};
      }else{
        $scope.component.service = $.extend(true,{},componentService.base_service);
        $scope.component.pod = $.extend(true,{},componentService.base_pod);
        if( $scope.component.service.spec.type == "NodePort"){
          $scope.component.service.spec.ports.push($.extend(true,{},componentService.nodeport));
        }else{
          $scope.component.service.spec.ports.push($.extend(true,{},componentService.clusterip));
        } 
      }
    }

    $scope.addPort = function(){
      if( $scope.component.service.spec.type == "NodePort"){
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.nodeport));
      }else{
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.clusterip));
      } 
    }

    $scope.removePort = function(index){
      $scope.component.service.spec.ports.splice(index,1);
    }

    $scope.changeServiceType = function(){
      if( $scope.component.service.spec.type == "NodePort"){
        $scope.component.service.spec.ports = [];
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.nodeport));
      }else{
        $scope.component.service.spec.ports = [];
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.clusterip));
      }
    }

    $scope.createImage = function(){
      $scope.toCreateImage = true;
      $scope.component.image_name = "";
      $scope.component.image_tag = "";
      $scope.component.image_setting = $.extend(true,{},componentService.imagesetting);
    }

    $scope.cancelCreateImage = function(){
      $scope.toCreateImage = false;
      $scope.component.image_setting = {};
    }
    
    $scope.switchMode = function(){
      $scope.jsonMode = !$scope.jsonMode;
      if($scope.jsonMode){
        componentIO.initFromEdit("input");
        componentIO.initFromEdit("output");
      }else{
        componentIO.initTreeEdit();
      }
    }

    // init component create page
    function init(){
      $scope.component = $.extend(true,{},componentService.component);

      // for runtime config tabs control
      $scope.runtimeTab = 1;

      // determine if to create image
      $scope.toCreateImage = false;

      // determine which editor to use for input output json
      $scope.jsonMode = false; 

      $scope.baseOrAdvanced();
    }

    init();
}]);