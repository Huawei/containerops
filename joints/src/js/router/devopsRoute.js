var devops = angular.module('DevOps',['ngRoute','ui.router','ngWebSocket']);

devops.run(['$rootScope', '$state', '$stateParams',
    function ($rootScope, $state, $stateParams) {
       $rootScope.$state = $state;
       $rootScope.$stateParams = $stateParams;
    }
  ]
).config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
      $urlRouterProvider.otherwise("/workflow");
      $stateProvider
        .state('home', {
          abstract: true,
          templateUrl: 'templates/home/main.html',
          controller: 'HomeController'
        })
        .state('workflow', {
          parent: 'home',
          url: '/workflow',
          views: {
            'main': {
              templateUrl: 'templates/workflow/main.html',
              controller: 'WorkflowController'
            }
          }
        })
        .state('component', {
          parent: 'home',
          url: '/component',
          views: {
            'main': {
              templateUrl: 'templates/component/main.html',
              controller: 'ComponentController'
            }
          }
        })
        .state('component.create', {
          parent: 'home',
          url: '/component/create',
          views: {
            'main': {
              templateUrl: 'templates/component/create.html',
              controller: 'CreateComponentController'
            }
          }
        })
        .state('component.detail', {
          parent: 'home',
          url: '/component/:id',
          views: {
            'main': {
              templateUrl: 'templates/component/detail.html',
              controller: 'ComponentDetailController'
            }
          }
        })
}]);
