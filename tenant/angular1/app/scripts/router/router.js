var auth = angular.module('Auth',['ui.router','ui.bootstrap']);

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
        .state('organization.create', {
          parent: 'organization',
          url: '/create',
          views: {
            'main@core': {
              templateUrl: 'templates/organization/create.html',
              controller: 'OrganizationCreateController'
            }
          }
        })
        .state('team', {
          parent: 'core',
          url: '/team',
          params:{'orgId':''},
          views: {
            'main': {
              templateUrl: 'templates/team/main.html',
              controller: 'TeamController'
            }
          }
        })
        .state('team.create', {
          parent: 'team',
          url: '/create',
          views: {
            'main@core': {
              templateUrl: 'templates/team/create.html',
              controller: 'TeamCreateController'
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
        .state('project.create', {
          parent: 'project',
          url: '/create',
          views: {
            'main@core': {
              templateUrl: 'templates/project/create.html',
              controller: 'ProjectCreateController'
            }
          }
        })
        .state('project.edit', {
          parent: 'project',
          url: '/edit',
          params: {name:null, id:null},
          views: {
            'main@core': {
              templateUrl: 'templates/project/create.html',
              controller: 'ProjectCreateController'
            }
          }
        })
        .state('application', {
          parent: 'core',
          url: '/application',
          params: {name:null, id:null},
          views: {
            'main': {
              templateUrl: 'templates/application/main.html',
              controller: 'ApplicationController'
            }
          }
        })
        .state('application.create', {
          parent: 'application',
          url: '/create',
          views: {
            'main@core': {
              templateUrl: 'templates/application/create.html',
              controller: 'ApplicationCreateController'
            }
          }
        })
        .state('application.edit', {
          parent: 'application',
          url: '/edit',
          params: {name:null, id:null},
          views: {
            'main@core': {
              templateUrl: 'templates/application/create.html',
              controller: 'ApplicationCreateController'
            }
          }
        })
        .state('module', {
          parent: 'core',
          url: '/module',
          params: {name:null, id:null},
          views: {
            'main': {
              templateUrl: 'templates/module/main.html',
              controller: 'ModuleController'
            }
          }
        })
        .state('module.create', {
          parent: 'module',
          url: '/create',
          views: {
            'main@core': {
              templateUrl: 'templates/module/create.html',
              controller: 'ModuleCreateController'
            }
          }
        })
        .state('module.edit', {
          parent: 'module',
          url: '/edit',
          params: {name:null, id:null},
          views: {
            'main@core': {
              templateUrl: 'templates/module/create.html',
              controller: 'ModuleCreateController'
            }
          }
        })
}]);
