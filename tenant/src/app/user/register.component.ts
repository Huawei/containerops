import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';

var md5 = require("blueimp-md5/js/md5")

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
		console.log(this.user);
		this.changeNav('login');
	}
}