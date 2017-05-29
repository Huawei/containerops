/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

define(['app'], function(app) {
    return app.run(['$rootScope', '$state', '$stateParams',
            function($rootScope, $state, $stateParams) {
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
            }
        ])
        .config(['$controllerProvider', '$provide', '$compileProvider', '$filterProvider',
            function($controllerProvider, $provide, $compileProvider, $filterProvider) {
                app.controllerProvider = $controllerProvider;
                app.provide = $provide;
                app.compileProvider = $compileProvider;
                app.filterProvider = $filterProvider;
            }
        ])
        .config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
            $urlRouterProvider.otherwise("/workflow");
            $stateProvider
                .state('home', {
                    abstract: true,
                    templateUrl: 'templates/home/main.html',
                    controller: 'HomeController',
                    resolve: {
                        loadCtrl_home: ["$q", function($q) {
                            var deferred = $q.defer();
                            require(["controllers/home/main"], function() {
                                deferred.resolve();
                            });
                            return deferred.promise;
                        }],
                    }
                })
                .state('workflow', {
                    parent: 'home',
                    url: '/workflow',
                    views: {
                        'main': {
                            templateUrl: 'templates/workflow/main.html',
                            controller: 'WorkflowController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/main"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('workflow.create', {
                    url: '/create',
                    views: {
                        'main@home': {
                            templateUrl: 'templates/workflow/create.html',
                            controller: 'WorkflowCreateController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/create"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('workflow.create.stage', {
                    url: '/stage/:id',
                    views: {
                        'element': {
                            templateUrl: 'templates/workflow/stageDetail.html',
                            controller: 'WorkflowElementController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/element"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('workflow.create.action', {
                    url: '/action/:id',
                    views: {
                        'element': {
                            templateUrl: 'templates/workflow/actionDetail.html',
                            controller: 'ActionDetailController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/actionDetail"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('workflow.create.addComponent', {
                    url: '/addComponent',
                    views: {
                        'element': {
                            templateUrl: 'templates/workflow/addComponent.html',
                            controller: 'AddComponentController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/addComponent"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('workflow.create.editComponent', {
                    url: '/editComponent/:id',
                    views: {
                        'element': {
                            templateUrl: 'templates/workflow/componentDetail.html',
                            controller: 'ComponentDetailController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/workflow/componentDetail"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('component', {
                    parent: 'home',
                    url: '/component',
                    views: {
                        'main': {
                            templateUrl: 'templates/component/main.html',
                            controller: 'ComponentController',
                            resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/component/main"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('component.create', {
                    url: '/create',
                    views: {
                        'main@home': {
                            templateUrl: 'templates/component/create.html',
                            controller: 'CreateComponentController',
                             resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/component/create"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
                .state('component.detail', {
                    url: '/:id',
                    views: {
                        'main@home': {
                            templateUrl: 'templates/component/detail.html',
                            controller: 'ComponentDetailController',
                             resolve: {
                                loadCtrl_workflow: ["$q", function($q) {
                                    var deferred = $q.defer();
                                    require(["controllers/component/detail"], function() {
                                        deferred.resolve();
                                    });
                                    return deferred.promise;
                                }],
                            }
                        }
                    }
                })
        }]);


})
