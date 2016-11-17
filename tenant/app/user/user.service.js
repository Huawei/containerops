"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require('@angular/core');
var platform_browser_1 = require('@angular/platform-browser');
var http_1 = require('@angular/http');
require('rxjs/add/operator/toPromise');
require('rxjs/add/operator/catch');
var UserService = (function () {
    function UserService(http, title) {
        this.http = http;
        this.title = title;
        this.headers = new http_1.Headers({ 'Content-Type': 'application/json' });
    }
    UserService.prototype.getBrowseList = function () {
        return this.http.get('json/browseList.json')
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.doLogin = function (info) {
        var params = JSON.stringify(info);
        return this.http.post('/web/v1/user/signin', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.signUp = function (info) {
        var params = JSON.stringify(info);
        return this.http.post('/web/v1/user', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.sendEmail = function (info) {
        var params = JSON.stringify(info);
        return this.http.post('/web/v1/user/forget', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.resetPwd = function (info) {
        var params = JSON.stringify(info);
        return this.http.post('/web/v1/user/forget/reset', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.getEmailList = function (user) {
        return this.http.get('/web/v1/user/' + user.username + '/emails')
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.addEmail = function (info) {
        var params = JSON.stringify(info);
        return this.http.put('/web/v1/user/' + info.username + '/email', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.verifyEmail = function (info, user) {
        info.username = user.username;
        var params = JSON.stringify(info);
        return this.http.put('/web/v1/user/' + user.username + '/email/' + info.id + '/send', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData)
            .catch(this.handleError);
    };
    UserService.prototype.delEmail = function (info, user) {
        var params = JSON.stringify(info);
        return this.http.delete('/web/v1/user/' + user.username + '/email/' + info.id, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.loginOut = function (user) {
        var params = JSON.stringify(user);
        return this.http.put('/web/v1/user/' + user.username + '/signout', params, { headers: this.headers })
            .toPromise()
            .then(this.dealData, this.dealError)
            .catch(this.handleError);
    };
    UserService.prototype.dealData = function (res) {
        var object = {
            code: res.status,
            data: res.json()
        };
        console.log(res);
        return object || {};
    };
    UserService.prototype.dealError = function (err) {
        var object = {
            code: err.status,
            data: err.json()
        };
        console.log(err);
        return object || {};
    };
    UserService.prototype.handleError = function (error) {
        console.log(error);
        // let errMsg = (error.message) ? error.message :
        //   error.status ? `${error.status} - ${error.statusText}` : 'Server error';
        // console.log(errMsg);
        var object = {
            code: error.status,
            data: error.json()
        };
        return Promise.reject(object);
    };
    UserService.prototype.changeTitle = function (val) {
        this.title.setTitle(val);
    };
    UserService = __decorate([
        core_1.Injectable(), 
        __metadata('design:paramtypes', [http_1.Http, platform_browser_1.Title])
    ], UserService);
    return UserService;
}());
exports.UserService = UserService;
