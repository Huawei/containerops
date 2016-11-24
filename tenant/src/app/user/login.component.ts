import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';

var md5 = require("blueimp-md5/js/md5");
var _ = require("underscore");

@Component({
  selector: 'login',
  templateUrl: '../../template/user/login.html'
})

export class LoginComponent implements OnInit { 
	user = {
		username: '',
		password: ''
	}

	constructor(private router: Router,private userService: UserService){

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
				alert("no such user, please sign up first.")
			}else{
				var users = JSON.parse(localStorage["users"]);
				var targetuser = _.find(users,function(item){
					return item.username == self.user.username && item.password == self.user.password;
				});
				if(_.isUndefined(targetuser)){
					alert("no such user, please sign up first.")
				}else{
					alert("sign in done.")
					this.changeNav('index');
				}
			}
			//fake end
		}catch(e){
			alert("failed to sign in.")
		}
	}
}