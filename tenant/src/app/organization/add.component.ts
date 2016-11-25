import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {OrgDataService} from './data.service';

@Component({
  selector: 'add-org',
  templateUrl: '../../template/organization/add.html'
})

export class OrgAddComponent implements OnInit { 
	org = {
		"name" : "",
		"desc" : ""
	};

	constructor(private router: Router, private orgDataService: OrgDataService){

	}

	ngOnInit(): void {

	}

	addOrg(): void {
		try{
			//fake, to be deleted
			this.orgDataService.addOrg(this.org);
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