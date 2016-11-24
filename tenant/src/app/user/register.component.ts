import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';

var md5 = require("blueimp-md5/js/md5");
var _ = require("underscore");

@Component({
  selector: 'register',
  templateUrl: '../../template/user/register.html'
})

export class RegisterComponent implements OnInit { 
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
			alert("sign up done.")
			//fake end
			this.changeNav('login');
		}catch(e){
			alert("failed to sign up.")
		}
	}
}