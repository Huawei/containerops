import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'list-repo',
  templateUrl: '../../template/repository/list.html'
})

export class RepoListComponent implements OnInit { 
	constructor(private router: Router){

	}

	ngOnInit(): void {
		
	}
}