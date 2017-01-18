devops.controller('ComponentDetailController', ['$scope','$stateParams', '$location', 'componentService', 'componentIO', 
  'componentCheck', 'notifyService', 'apiService', 'loading',
  function($scope,$stateParams,$location,componentService,componentIO,componentCheck,notifyService,apiService,loading) {   

    // tabs control
    $scope.selectTab = function(index){
      var result = checkTabs();
      if(index == 2){
        initScriptEditor();
        showEventScript();
      }
      // else if(index == 4 && !result){
      //   return;
      // }

      $scope.tab = index;
    }

    function checkTabs(){
      // var result = true;
      $scope.tabStatus.runtime = componentCheck.tabcheck.runtime($scope.toCreateImage);
      // result = result && $scope.tabStatus.runtime;
      if($scope.toCreateImage){
        $scope.tabStatus.editshell = componentCheck.tabcheck.editshell();
        $scope.tabStatus.buildimage = componentCheck.tabcheck.buildimage();
        // result = result && $scope.tabStatus.editshell && $scope.tabStatus.buildimage;
      }
      // return result;
    }

    // runtime tab all functions below
    $scope.changeRuntimeTab = function(index){
      $scope.runtimeTab = index;
      if($scope.runtimeTab == 2){
        componentIO.init($scope.component);
      }
    }

    $scope.baseOrAdvanced = function(){
      if($scope.component.use_advanced){
        $scope.component.service = {};
        $scope.component.pod = {};

        setTimeout(function(){
          $("#serviceCodeEditor").val(JSON.stringify($scope.component.service,null,2));
          $("#serviceCodeEditor").on("blur",function(){
              var result = toJsonYaml("service");
              if(result){
                  $scope.component.service = result;
              }    
          });

          $("#podCodeEditor").val(JSON.stringify($scope.component.pod,null,2));
          $("#podCodeEditor").on("blur",function(){
              var result = toJsonYaml("pod");
              if(result){
                  $scope.component.pod = result;
              }    
          });
        },500);
      }else{
        $scope.component.service = $.extend(true,{},componentService.metadata.base_service);
        $scope.component.pod = $.extend(true,{},componentService.metadata.base_pod);
        if( $scope.component.service.spec.type == "NodePort"){
          $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.nodeport));
        }else{
          $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.clusterip));
        } 
        $scope.serviceType = $scope.component.service.spec.type;
      }
    }

    $scope.addPort = function(){
      if( $scope.component.service.spec.type == "NodePort"){
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.nodeport));
      }else{
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.clusterip));
      } 
    }

    $scope.removePort = function(index){
      $scope.component.service.spec.ports.splice(index,1);
    }

    $scope.changeServiceType = function(){
      if( $scope.component.service.spec.type == "NodePort"){
        $scope.component.service.spec.ports = [];
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.nodeport));
      }else{
        $scope.component.service.spec.ports = [];
        $scope.component.service.spec.ports.push($.extend(true,{},componentService.metadata.clusterip));
      }
      $scope.serviceType = $scope.component.service.spec.type;
    }

    $scope.createImage = function(){
      $scope.toCreateImage = true;
      $scope.component.image_name = "";
      $scope.component.image_tag = "";
      $scope.component.image_setting = $.extend(true,{},componentService.metadata.imagesetting);
    }

    $scope.cancelCreateImage = function(){
      $scope.toCreateImage = false;
      $scope.component.image_name = "";
      $scope.component.image_tag = "";
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

    $scope.addEnv = function(){
       $scope.component.env.push($.extend(true,{},componentService.metadata.env));
    }

    $scope.removeEnv = function(index){
      $scope.component.env.splice(index,1);
    }

    function toJsonYaml(type){
        var value,result;
        if(type == "service"){
            value = $("#serviceCodeEditor").val();
        }else if(type == "pod"){
            value = $("#podCodeEditor").val();
        }

        try{
            result = JSON.parse(value);
        }catch(e){
           try{
            result = jsyaml.safeLoad(value);
           }catch(e){
            notifyService.notify("Your advanced " + type + " setting is not a legal json or yaml.","error");
            result = false;
           }
        }
        if(!result){
           notifyService.notify("Your advanced " + type + " setting is not a legal json or yaml.","error");
        }
        return result;
    }

    // shell edit page functions
    $scope.selectEvent = function(index){
      $scope.selectedEvent = index;
      showEventScript();
    }

    function initScriptEditor(){
      $scope.scriptEditor = ace.edit("scriptEditor");
      $scope.scriptEditor.setTheme("ace/theme/dawn");
      $scope.scriptEditor.getSession().setMode("ace/mode/golang");
      $scope.scriptEditor.on("blur",function(){
        setEventScript();
      })
    }

    function showEventScript(){
      switch($scope.selectedEvent){
        case 1: 
            $scope.scriptEditor.setValue($scope.component.image_setting.events.component_start); 
            break;
        case 2: 
            $scope.scriptEditor.setValue($scope.component.image_setting.events.component_result); 
            break;
        case 3: 
            $scope.scriptEditor.setValue($scope.component.image_setting.events.component_stop); 
            break;
      }

      var session = $scope.scriptEditor.getSession();
      var row = session.getLength();
      var col = session.getLine(row-1).length;
      $scope.scriptEditor.gotoLine(row,col);
      $scope.scriptEditor.focus();
    }

    function setEventScript(){
      switch($scope.selectedEvent){
        case 1: 
            $scope.component.image_setting.events.component_start = $scope.scriptEditor.getValue(); 
            break;
        case 2: 
            $scope.component.image_setting.events.component_result = $scope.scriptEditor.getValue(); 
            break;
        case 3: 
            $scope.component.image_setting.events.component_stop = $scope.scriptEditor.getValue(); 
            break;
      }
    }

    // build push page functions
    $scope.changeBaseImageType = function(){
      if($scope.component.image_setting.from.type == "url"){
        delete $scope.component.image_setting.from.name;
        delete $scope.component.image_setting.from.tag;
        delete $scope.component.image_setting.from.dockerfile;
        $scope.component.image_setting.from["url"] = "";
      }else if($scope.component.image_setting.from.type == "dockerfile"){
        delete $scope.component.image_setting.from.name;
        delete $scope.component.image_setting.from.tag;
        delete $scope.component.image_setting.from.url;
        $scope.component.image_setting.from["dockerfile"] = "";
      }else if($scope.component.image_setting.from.type == "name"){
        delete $scope.component.image_setting.from.url;
        delete $scope.component.image_setting.from.dockerfile;
        $scope.component.image_setting.from["name"] = "";
        $scope.component.image_setting.from["tag"] = "";
      }
    }

    // save component
    $scope.saveComponent = function(){
      var result = componentCheck.go($scope.toCreateImage);
      if(result){
        var promise = componentService.updateComponent($scope.component);
        promise.done(function(data){
            apiService.successToCall(data);
            $location.path("/component/" + data.component.id);
        });
        promise.fail(function(xhr,status,error){
            apiService.failToCall(xhr.responseJSON);
        });
      }
    }

    // back to list
    $scope.backToList = function(){
      $location.path("/component");
    }

    // get data
  	function getComponent(id){
  		var promise = componentService.getComponent(id);
  		promise.done(function(data){
  			loading.hide();
  			$scope.component = data.component;

        // determine if to create image
        $scope.toCreateImage = _.isEmpty($scope.component.image_name);

  			// for top tabs control
  			$scope.tab = 1;
  			$scope.tabStatus = {
  				"runtime" : true,
  				"editshell" : $scope.toCreateImage,
  				"buildimage" : $scope.toCreateImage
  			}

		    // for runtime config tabs control
		    $scope.runtimeTab = 1;

		    // determine which editor to use for input output json
		    $scope.jsonMode = false;

		    // for event selection
		    $scope.selectedEvent = 1;

		    // init service pod
		    $scope.serviceType = $scope.component.service.spec.type;

		    // init check
		    componentCheck.init($scope.component);

		    $scope.dataReady = true;
		    $scope.$apply();
  		});
  		promise.fail(function(xhr,status,error){
  			apiService.failToCall(xhr.responseJSON);
  		}); 
  	}

  	// init component detail page
  	function init(){
  		$scope.dataReady = false;
  		getComponent($stateParams.id);
  	}
    
    init();
}]);
