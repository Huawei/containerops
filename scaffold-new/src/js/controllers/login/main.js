define(['login','services/login/main'], function (app) {
    'use strict';
     app.controllerProvider.register('LoginController', ['$scope','$uibModal','$window','$cookies','LoginService','ResponseService',function ($scope,$uibModal,$window,$cookies,LoginService,ResponseService) {
    	  $scope.user = {"email":"","password":""};
        $scope.newuser = {"newuser":"","email":"","password":""};
    	  $scope.login = function(){
           LoginService.doLogin($scope.user).then(function(data){
                if(data.code == 200){
                    var expireDate = new Date();
                    expireDate.setDate(expireDate.getDate() + 2);
                    // Setting a cookie
                
                    $cookies.put('email',$scope.user.email,{'expires': expireDate});
                    $window.location = "index.html";    
                }else{
                  var error = {
                    msg : data.msg,
                    code : data.code
                  }
                  ResponseService.errorResponse(error);  
                }
                           
             },
              function(error){
				         ResponseService.errorResponse(error);		    
			       })
    	  };

        $scope.signup = function(){
           LoginService.doSignup($scope.newuser).then(function(data){
                if(data.code == 200){
                  alert("Sign up successfully, please login.")
                    $window.location = "/";    
                }else{
                  var error = {
                    msg : data.msg,
                    code : data.code
                  }
                  ResponseService.errorResponse(error);  
                }
                           
             },
              function(error){
                 ResponseService.errorResponse(error);        
             })
        };

     }]);
});