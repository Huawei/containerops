import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';
import { NotifyService } from '../common/notify.service';

var md5 = require("blueimp-md5/js/md5");
var _ = require("underscore");

@Component({
  selector: 'login',
  template: require('../../template/user/login.html')
})

export class LoginComponent implements OnInit { 
	user = {
		username: '',
		password: ''
	}

	constructor(private router: Router,
				private userService: UserService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
		
	}

	changeNav(val){
    	this.router.navigate([val]);
	}

	login() {
		try{
			//fake, to be deleted
			var self = this;
			if(_.isUndefined(localStorage["users"])){
				this.notifyService.notify("No such user, please sign up first.","info");
			}else{
				var users = JSON.parse(localStorage["users"]);
				var targetuser = _.find(users,function(item){
					return item.username == self.user.username && item.password == self.user.password;
				});
				if(_.isUndefined(targetuser)){
					this.notifyService.notify("No such user or password is incorrect.","info");
				}else{
					this.notifyService.notify("Welcome. " + self.user.username , "success");
					sessionStorage["currentUser"] = self.user.username;
					this.changeNav('index');
				}
			}
			//fake end
		}catch(e){
			this.notifyService.notify("Failed to Sign in.", "error");
		}
	}
}