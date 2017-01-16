auth.factory('OrganizationService', ['$http', '$q',  function($http, $q) {
    return {
        get: function() {
            var deferred = $q.defer();
            var url = "/organization";
            var request = {
                "url": url,
                "dataType": "json",
                "method": "GET"
            }
            $http(request).then(function(data) {
                deferred.resolve(data);
            }, function(error) {
                deferred.reject(error);
            });
            return deferred.promise;
        },
        save: function(data) {
            var deferred = $q.defer();
            var url = "/organization";
            var request = {
                "url": url,
                "dataType": "json",
                "method": "POST",
                "data": angular.toJson(data)
            }
            $http(request).then(function(data) {
                deferred.resolve(data);
            }, function(error) {
                deferred.reject(error);
            });
            return deferred.promise;
        }
       
    }

}]);
