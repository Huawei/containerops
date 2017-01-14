devops.controller('HomeController', ['$scope','$location',function($scope,$location) {  
    
	  $scope.navigators = [        
    		{
          "name" : "Workflow",
          "href" :"/workflow",
          "ngclass" : "active",
          "icon" : "fa fa-desktop"
        },
    		{
          "name" : "Component",
          "href" : "/component",
          "ngclass" : "",
          "icon": "fa fa-cube"
        },
    		{
          "name" : "History",
          "href" : "/history",
          "ngclass" : "",
          "icon" : "fa fa-history"
        },
    		{
          "name" : "System Setting",
          "href" : "/setting",
          "ngclass" : "",
          "icon" : "fa fa-cog"
        }
  	];

  	$scope.chooseNav = function(name){
  		_.each($scope.navigators,function(item){
  			if(item.name != name){
  				item.ngclass = "";
  			}else{
  				item.ngclass = "active";
  			}
  		})
  	}
  	
  	function initializeNav(){
  		var path = $location.path();
  		_.each($scope.navigators,function(item){
  			if(path != ""){
  				if(path.indexOf(item.href) == 0){
            item.ngclass = "active";
          }else{
            item.ngclass = "";
          }
  			}
  		})
  	}
  	
  	initializeNav();

}]);