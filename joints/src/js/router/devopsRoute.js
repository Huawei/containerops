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
