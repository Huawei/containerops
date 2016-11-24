import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

var _ = require("underscore");

@Component({
  selector: 'add-org',
  templateUrl: '../../template/organization/add.html'
})

export class OrgAddComponent implements OnInit { 
	org = {
		"name" : "",
		"desc" : ""
	};

	constructor(private router: Router){

	}

	ngOnInit(): void {

	}

	addOrg(): void {
		try{
			//fake, to be deleted
			var orgs;
			if(_.isUndefined(localStorage["orgs"])){
				orgs = [];
			}else{
				orgs = JSON.parse(localStorage["orgs"])
			}
			orgs.push(this.org);
			localStorage["orgs"] = JSON.stringify(orgs);
			alert("org "+ this.org.name + " added.");

			this.router.navigate(["organization"]);
			//fake end
		}catch(e){
			alert("failed to add org.")
		}
	}

	cancelAdd(): void{
		this.router.navigate(["organization"]);
	}	
}