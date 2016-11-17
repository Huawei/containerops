import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { UserService } from './user.service';
var md5 = require("blueimp-md5/js/md5")

@Component({
  selector: 'register',
  templateUrl: '../../template/user/register.html'
})

export class RegisterComponent implements OnInit { 
	errorMsg: string;
	isTips = {
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
  }
	user = {
		username: '',
		email: '',
		password: ''
	}
	active = '';
	browseList = [];
	hover = '';
	salt = '';

	constructor(
		private router: Router,
	 	private userService: UserService){
			this.changeTitle('- register')
	}

	ngOnInit(): void {
		var salt = document.getElementsByTagName('meta')['salt'].getAttribute('content')
		this.salt = salt;
	}

	changeTitle(val) {
		var title = (document.getElementsByTagName('title')[0].innerHTML) ? (document.getElementsByTagName('title')[0].innerHTML).split('-')[0] + val : val;
		this.userService.changeTitle(title)
	}
	
  activeHover(index){
  	this.hover = index;
  }

  changeNav(val){
		// this.active = val
    this.router.navigate([val]);
	}

	signUp() {
		console.log(this.user)
		var user = this.user;
		// var pwdReg=/(?![0-9a-z]+$)(?![a-zA-Z]+$){8,}/
		// console.log(/(?![0-9a-z]+$)(?![a-zA-Z]+$){8,}/.test(user.password))
		var nameReg = /[0-9A-Za-z]{1,}/.test(user.username)
		var password = user.password

		var pwdReg = password && (password.length > 8) && (password.indexOf(''+user.username)===-1) && (/[0-9]/g.test(password)) && (/[A-Z]/g.test(password))
		if(nameReg&&pwdReg&&user.email&&user.email.indexOf('@')!==-1){
			var data = {
				username: user.username,
				email: user.email,
				password: md5(this.salt + user.password)
			}
			this.userService.signUp(data)
      .then(res => { 
      	if(res.code === 201){
      		this.router.navigate(['repositories'])
      		sessionStorage.setItem("username", user.username)
      	}else if(400 <= res.code && res.code < 500){
      		this.tips('otherError',true)
      		this.isTips.otherText = res.data.message
      	}
      },error => {
      	console.log(error)
      	if(400 <= error.code && error.code < 500){
      		this.tips('otherError',true)
      		this.isTips.otherText = error.data.message
      	}
      });
		}else if(!user.username){
			this.tips('username',true)
		}else if(!(user.username&&nameReg)){
			console.log('have username')
			this.tips('isUsernameRight',true)
		}else if(!user.email){
			console.log('no email')
			// this.isTips.emailInfo = 'email is required'
			this.tips('email',true)
		}else if(!(user.email.indexOf('@')!==-1)){
			console.log('have email')
			// this.isTips.emailInfo = 'email is invalid'
			this.tips('emailError',true)
		}else if(user.password){
			console.log('have password')
			if(password.length < 8){
				this.tips('password',true)
				this.isTips.passwordText = 'password at least eight characters'
			}else if(user.password.indexOf(''+user.username)!==-1){
				this.tips('password',true)
				this.isTips.passwordText = 'password cannot contain the user name'
			}else if(!(/[0-9]/g.test(password))){
				this.tips('password',true)
				this.isTips.passwordText = 'password must contain a number'
			}else if(!(/[A-Z]/g.test(password))){
				this.tips('password',true)
				this.isTips.passwordText = 'password must contain a capital letter'
			}
		}else if(!user.password){
			console.log('no password')
			this.tips('password',true)
			this.isTips.passwordText = 'password is required'
		}
	}

	tips(name,val){
		this.isTips[name] = val;
		if(val){
			setTimeout(function(){
				this.tips(name,false)
			}.bind(this),3000)
		}
	}
}