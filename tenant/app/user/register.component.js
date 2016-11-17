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
var router_1 = require('@angular/router');
var user_service_1 = require('./user.service');
var md5 = require("blueimp-md5/js/md5");
var RegisterComponent = (function () {
    function RegisterComponent(router, userService) {
        this.router = router;
        this.userService = userService;
        this.isTips = {
            username: false,
            isUsernameRight: false,
            email: false,
            // emailInfo: '',
            emailError: false,
            emailInfo: false,
            password: false,
            pwdError: false,
            otherError: false,
            otherText: '',
            passwordText: ''
        };
        this.user = {
            username: '',
            email: '',
            password: ''
        };
        this.active = '';
        this.browseList = [];
        this.hover = '';
        this.salt = '';
        this.changeTitle('- register');
    }
    RegisterComponent.prototype.ngOnInit = function () {
        var salt = document.getElementsByTagName('meta')['salt'].getAttribute('content');
        this.salt = salt;
    };
    RegisterComponent.prototype.changeTitle = function (val) {
        var title = (document.getElementsByTagName('title')[0].innerHTML) ? (document.getElementsByTagName('title')[0].innerHTML).split('-')[0] + val : val;
        this.userService.changeTitle(title);
    };
    RegisterComponent.prototype.activeHover = function (index) {
        this.hover = index;
    };
    RegisterComponent.prototype.changeNav = function (val) {
        // this.active = val
        this.router.navigate([val]);
    };
    RegisterComponent.prototype.signUp = function () {
        var _this = this;
        console.log(this.user);
        var user = this.user;
        // var pwdReg=/(?![0-9a-z]+$)(?![a-zA-Z]+$){8,}/
        // console.log(/(?![0-9a-z]+$)(?![a-zA-Z]+$){8,}/.test(user.password))
        var nameReg = /[0-9A-Za-z]{1,}/.test(user.username);
        var password = user.password;
        var pwdReg = password && (password.length > 8) && (password.indexOf('' + user.username) === -1) && (/[0-9]/g.test(password)) && (/[A-Z]/g.test(password));
        if (nameReg && pwdReg && user.email && user.email.indexOf('@') !== -1) {
            var data = {
                username: user.username,
                email: user.email,
                password: md5(this.salt + user.password)
            };
            this.userService.signUp(data)
                .then(function (res) {
                if (res.code === 201) {
                    _this.router.navigate(['repositories']);
                    sessionStorage.setItem("username", user.username);
                }
                else if (400 <= res.code && res.code < 500) {
                    _this.tips('otherError', true);
                    _this.isTips.otherText = res.data.message;
                }
            }, function (error) {
                console.log(error);
                if (400 <= error.code && error.code < 500) {
                    _this.tips('otherError', true);
                    _this.isTips.otherText = error.data.message;
                }
            });
        }
        else if (!user.username) {
            this.tips('username', true);
        }
        else if (!(user.username && nameReg)) {
            console.log('have username');
            this.tips('isUsernameRight', true);
        }
        else if (!user.email) {
            console.log('no email');
            // this.isTips.emailInfo = 'email is required'
            this.tips('email', true);
        }
        else if (!(user.email.indexOf('@') !== -1)) {
            console.log('have email');
            // this.isTips.emailInfo = 'email is invalid'
            this.tips('emailError', true);
        }
        else if (user.password) {
            console.log('have password');
            if (password.length < 8) {
                this.tips('password', true);
                this.isTips.passwordText = 'password at least eight characters';
            }
            else if (user.password.indexOf('' + user.username) !== -1) {
                this.tips('password', true);
                this.isTips.passwordText = 'password cannot contain the user name';
            }
            else if (!(/[0-9]/g.test(password))) {
                this.tips('password', true);
                this.isTips.passwordText = 'password must contain a number';
            }
            else if (!(/[A-Z]/g.test(password))) {
                this.tips('password', true);
                this.isTips.passwordText = 'password must contain a capital letter';
            }
        }
        else if (!user.password) {
            console.log('no password');
            this.tips('password', true);
            this.isTips.passwordText = 'password is required';
        }
    };
    RegisterComponent.prototype.tips = function (name, val) {
        this.isTips[name] = val;
        if (val) {
            setTimeout(function () {
                this.tips(name, false);
            }.bind(this), 3000);
        }
    };
    RegisterComponent = __decorate([
        core_1.Component({
            selector: 'register',
            templateUrl: '../../template/user/register.html'
        }), 
        __metadata('design:paramtypes', [router_1.Router, user_service_1.UserService])
    ], RegisterComponent);
    return RegisterComponent;
}());
exports.RegisterComponent = RegisterComponent;
