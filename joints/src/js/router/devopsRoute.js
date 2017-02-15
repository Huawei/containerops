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
                    parent: 'home',
                    url: '/component/create',
                    views: {
                        'main': {
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
                    parent: 'home',
                    url: '/component/:id',
                    views: {
                        'main': {
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
