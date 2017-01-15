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
        },
        terminateCluster: function(clusterid) {
            var deferred = $q.defer();
            var url = "/cluster/" + clusterid;
            var request = {
                "url": url,
                "method": "DELETE",
                "headers": {
                    "Content-Type": "application/json;charset=utf-8"
                }
            };
            $http(request).success(function(data) {
                deferred.resolve(data);
            }).error(function(error) {
                deferred.reject(error);
            });
            return deferred.promise;
        },
        validateClusterForUser: function(clustername) {
            var deferred = $q.defer();
            var url = '/clusterValidate';
            var request = {
                "url": url,
                "dataType": "json",
                "method": "GET",
                "params": {
                    "clustername": clustername,
                }
            };
            $http(request).success(function(data) {
                deferred.resolve(data);
            }).error(function(error) {
                deferred.reject(error);
            });
            return deferred.promise;
        }
        // deleteConfirm: function($scope) {
        //     $uibModal.open({
        //         templateUrl: 'templates/node/confirm.html',
        //         controller: 'ConfirmController',
        //         size: 'sm',
        //         backdrop: 'static',
        //         resolve: {
        //             model: function() {
        //                 return $scope.confirm;
        //             }
        //         }
        //     });
        // }
    }

}]);
