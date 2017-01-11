define(['login'], function (app) {
    'use strict';
     app.provide.factory('LoginService', ['$http','$q',function ($http,$q) {
	    	return {
                doLogin : function(user){
					var deferred = $q.defer();
					var url = "/user/login";		
					var request = {
						"url": url,
						"dataType": "json",
						"method": "POST",
						"data":angular.toJson(user)
					}
						
					$http(request).success(function(data){
						deferred.resolve(data);
					}).error(function(error){
						deferred.reject(error);
					});
					return deferred.promise;
					
			    },

			    doSignup : function(user){
					var deferred = $q.defer();
					var url = "/user/signup";		
					var request = {
						"url": url,
						"dataType": "json",
						"method": "POST",
						"data":angular.toJson(user)
					}
						
					$http(request).success(function(data){
						deferred.resolve(data);
					}).error(function(error){
						deferred.reject(error);
					});
					return deferred.promise;
					
			    }
			  
	    	}
	    	
     }]);
});