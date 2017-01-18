auth.factory('TeamService', ['$http', '$q',  function($http, $q) {
    return {
        get: function() {
            var deferred = $q.defer();
            var url = "/team";
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
            var url = "/team";
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
        },
        saveMember: function(teamid, user){
            var deferred = $q.defer();
            var url = "/team/"+teamid+"/saveMember";
            var request = {
                "url": url,
                "dataType": "json",
                "method": "POST",
                "data": angular.toJson(user)
            }
            $http(request).then(function(data) {
                deferred.resolve(data);
            }, function(error) {
                deferred.reject(error);
            });
            return deferred.promise;
        },
        getMembers: function(teamid){
            var deferred = $q.defer();
            var url = "/team/"+teamid+"/members";
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
        removeMember: function(teamid, member){
            var deferred = $q.defer();
            var url = "/team/"+teamid+"/members";
            var request = {
                "url": url,
                "dataType": "json",
                "method": "DELETE"
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
