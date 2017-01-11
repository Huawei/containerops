var login = angular.module('login', ['ngRoute']);

login.config(function($routeProvider, $locationProvider) {
	$routeProvider
		.when('/', {
			templateUrl: 'templates/user/login.html',
			controller: 'LoginController'
		})
		.when('/signup', {
			templateUrl: 'templates/user/register.html',
			controller: 'RegisterController'
		})
		.otherwise({
			redirectTo: '/'
		});
});

