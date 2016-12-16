import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { RepoDataService } from './data.service';
import { NotifyService } from '../common/notify.service';

@Component({
  selector: 'add-repo',
  template: require('../../template/repository/add.html')
})

export class RepoAddComponent implements OnInit { 
	repo = {
		"name" : "",
		"url" : "",
		"imagetag" : ""
	};

	constructor(private router: Router, 
				private repoDataService: RepoDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {

	}

	addRepo(): void {
		try{
			//fake, to be deleted
			this.repoDataService.addRepo(this.repo);
			this.notifyService.notify("Add repository '" + this.repo.name + "' successfully.","success");

			this.router.navigate(["repository"]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add repository.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(["repository"]);
	}	
}