import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'list-org',
  templateUrl: '../../template/organization/list.html'
})

export class OrgListComponent implements OnInit { 
	constructor(private router: Router){

	}

	ngOnInit(): void {
		
	}
}