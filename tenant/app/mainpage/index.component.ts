import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

var _ = require("underscore");

@Component({
  selector: 'app-index',
  templateUrl: '../../template/mainpage/index.html'
})

export class IndexComponent implements OnInit { 
	selectedNav;

	constructor(private router: Router){
		if(_.isUndefined(this.selectedNav)){
			this.selectedNav = "org";
		}
		
		this.router.navigate(['organizations']);
	}

	ngOnInit(): void {}

	selectNav(value){
		this.selectedNav = value;
	}
}