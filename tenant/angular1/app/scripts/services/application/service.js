/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */

function applicationService($http,$q){
  return {
    getList: function(params){
      var deferred = $q.defer();
      var url = "/applications";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    },
    saveBaseInfo: function(params){
      var deferred = $q.defer();
      var url = "/projects";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    },
    saveSetting: function(params){
      var deferred = $q.defer();
      var url = "/save";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    },
    getProjectList: function(params){
      var deferred = $q.defer();
      var url = "/orgs";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    },
    getTeamList: function(params){
      var deferred = $q.defer();
      var url = "/orgs";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    },
    getEditInfo: function(params){
      var deferred = $q.defer();
      var url = "/orgs";
      var request = {
          "url": url,
          "method": "GET"
      }
      $http(request).then(function(data) {
          deferred.resolve(data);
      }, function(error) {
          deferred.reject(error);
      });
      return deferred.promise;
    }
  }
}
   
auth.factory('applicationService', ['$http','$q',applicationService]);
