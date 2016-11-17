import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';
var md5 = require("blueimp-md5/js/md5")
// var $ = require("jquery/dist")


@Component({
  selector: 'login',
  templateUrl: '../../template/user/login.html'
})

export class LoginComponent implements OnInit { 
	// errorMsg: string;
	// isTips = false;
	// errorText = '';
	// user = {
	// 	username: '',
	// 	password: ''
	// }
	// active = '';
	// browseList = [];
	// hover = '';
	// salt = '';

	// constructor(
	// 	private router: Router,
	//  	private userService: UserService){
	// 	this.changeTitle('- login')
	// }

	// ngOnInit(): void {
	// 	var salt = document.getElementsByTagName('meta')['salt'].getAttribute('content')
	// 	this.salt = salt;
	// }

	// changeTitle(val) {
	// 	var title = (document.getElementsByTagName('title')[0].innerHTML) ? (document.getElementsByTagName('title')[0].innerHTML).split('-')[0] + val : val;
	// 	this.userService.changeTitle(title)
	// }
	
 //  activeHover(index){
 //  	this.hover = index;
 //  }

	// changeNav(val){
	// 	// this.active = val
 //    this.router.navigate([val]);
	// }

	// login() {
	// 	console.log(this.user)
	// 	var user = this.user;
	// 	if(user.username&&user.password){
	// 		var data = {
	// 			username: user.username,
	// 			password: md5(this.salt + user.password)
	// 		}
	// 		this.userService.doLogin(data)
 //      .then(res => { 
 //      	if(res.code === 200){
 //      		sessionStorage.setItem("username", res.data.username)
 //      		this.router.navigate(['repositories'])
 //      	}else if(400 <= res.code && res.code < 500){
 //      		this.tips(true)
 //      		this.errorText = res.data.message
 //      	}
 //      },error => {
 //      	if(400 <= error.code && error.code < 500){
 //      		this.tips(true)
 //      		this.errorText = error.data.message
 //      	}
 //      });
	// 	}
	// 	// else if(!user.username){
	// 	// 	this.tips.username = true;
	// 	// }else if(!user.password){
	// 	// 	this.tips.password = true;
	// 	// }
		
 //    // this.router.navigate(['repositories']);
	// }

	// // toLogin(envet){
	// // 	console.log(event)
	// // }
	// tips(val){
	// 	this.isTips = val
	// 	if(val){
	// 		setTimeout(function(){
	// 			this.tips(false)
	// 		}.bind(this),4000)
	// 	}
	// }
}