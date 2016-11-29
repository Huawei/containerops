import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {OrgDataService} from './data.service';
import { NotifyService } from '../common/notify.service';

@Component({
  selector: 'add-org',
  template: require('../../template/organization/add.html')
})

export class OrgAddComponent implements OnInit { 
	org = {
		"name" : "",
		"desc" : ""
	};

	constructor(private router: Router, 
				private orgDataService: OrgDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {

	}

	addOrg(): void {
		try{
			//fake, to be deleted
			this.orgDataService.addOrg(this.org);
			this.notifyService.notify("Add organization '" + this.org.name + "' successfully.","success");

			this.router.navigate(["organization"]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add organization.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(["organization"]);
	}	
}