var auth = angular.module('Auth',['ui.router']);

auth.run(['$rootScope', '$state', '$stateParams',
    function ($rootScope, $state, $stateParams) {
       $rootScope.$state = $state;
       $rootScope.$stateParams = $stateParams;
    }
  ]
).config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
      $urlRouterProvider.otherwise("/organization");
      $stateProvider
        .state('core', {
          abstract: true,
          templateUrl: 'templates/core/navigation.html',
          controller: 'CoreController'
        })
        .state('organization', {
          parent: 'core',
          url: '/organization',
          views: {
            'main': {
              templateUrl: 'templates/organization/main.html',
              controller: 'OrganizationController'
            }
          }
        })
        .state('team', {
          parent: 'core',
          url: '/team',
          views: {
            'main': {
              templateUrl: 'templates/team/main.html',
              controller: 'TeamController'
            }
          }
        })
        .state('project', {
          parent: 'core',
          url: '/project',
          views: {
            'main': {
              templateUrl: 'templates/project/main.html',
              controller: 'ProjectController'
            }
          }
        })
        .state('project.edit', {
          parent: 'core',
          url: '/projectEdit',
          views: {
            'main': {
              templateUrl: 'templates/project/edit.html',
              controller: 'ProjectEditController'
            }
          }
        })
        .state('application', {
          parent: 'core',
          url: '/application',
          views: {
            'main': {
              templateUrl: 'templates/application/main.html',
              controller: 'ApplicationController'
            }
          }
        })
        .state('module', {
          parent: 'core',
          url: '/module',
          views: {
            'main': {
              templateUrl: 'templates/module/main.html',
              controller: 'ModuleController'
            }
          }
        })
}]);
