import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';
import { NotifyService } from '../common/notify.service';

var md5 = require("blueimp-md5/js/md5");
var _ = require("underscore");

@Component({
  selector: 'register',
  template: require('../../template/user/register.html')
})

export class RegisterComponent implements OnInit { 
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

	signUp() {
		try{
			//fake, to be deleted
			var users = localStorage["users"];
			if(_.isUndefined(users)){
				users = [];
				users.push(this.user);
				localStorage["users"] = JSON.stringify(users);
			}else{
				users = JSON.parse(users);
				users.push(this.user);
				localStorage["users"] = JSON.stringify(users);
			}
			this.notifyService.notify("Sign up successfully. please login with your username.","success");
			//fake end
			this.changeNav('login');
		}catch(e){
			this.notifyService.notify("Fail to sign up.","error");
		}
	}
}