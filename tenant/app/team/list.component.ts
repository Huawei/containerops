import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'list-team',
  templateUrl: '../../template/team/list.html'
})

export class TeamListComponent implements OnInit { 
	constructor(private router: Router){

	}

	ngOnInit(): void {
		
	}
}