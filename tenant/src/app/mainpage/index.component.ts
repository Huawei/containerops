import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
var _ = require("underscore");
// declare var _: any;

@Component({
  selector: 'app-index',
  template: require('../../template/mainpage/index.html')

})

export class IndexComponent implements OnInit { 
	currentUser;
	selectedNav;

	constructor(private router: Router){
		this.currentUser = sessionStorage["currentUser"]; 

		if(_.isUndefined(this.selectedNav)){
			this.selectedNav = "org";
		}
		
		this.router.navigate(['organization']);
	}

	ngOnInit(): void {}

	selectNav(value){
		this.selectedNav = value;
	}

	toggleSideBar(event): void{
		var target = $(event.currentTarget);
	    if(target.hasClass("sidebar-close")){
	        target.removeClass("sidebar-close").addClass("sidebar-open");
	        target.removeClass("fa-chevron-circle-left").addClass("fa-chevron-circle-right");
	        $("body").removeClass("nav-static").addClass("nav-collapsed");
	    }else if(target.hasClass("sidebar-open")){
	        target.removeClass("sidebar-open").addClass("sidebar-close");
	        target.removeClass("fa-chevron-circle-right").addClass("fa-chevron-circle-left");
	        $("body").removeClass("nav-collapsed").addClass("nav-static");
	    }
	}
}