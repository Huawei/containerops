import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';

var md5 = require("blueimp-md5/js/md5")

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
		console.log(this.user);
		this.changeNav('index');
	}
}