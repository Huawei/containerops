import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { TeamDataService } from './data.service';
import { UserService } from '../user/user.service';
import { NotifyService } from '../common/notify.service';

@Component({
  selector: 'team-add-member',
  templateUrl: '../../template/team/addMember.html'
})

export class TeamAddMemberComponent implements OnInit { 
    team;
	users;
    selectedUserId;
	constructor(private router: Router,
				private route: ActivatedRoute,
				private teamDataService: TeamDataService,
				private userService: UserService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
		this.getUsers();
		this.getTeam();
	}


	getTeam() {
		this.route.params.subscribe(params => {
			let id = +params['id'];
			this.team = this.teamDataService.getTeam(id);
		})
	}
	getUsers(): void {
		this.users = this.userService.getUsers();
		this.selectedUserId = this.users[0].username;
	}

	addMember(): void {
		try{
			//fake, to be deleted
			this.userService.updateUser(this.selectedUserId, this.team);
			this.notifyService.notify("Add member to '" + this.team.name + "' successfully.","success");

			this.router.navigate(['/team', this.team.id]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add member.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(['/team', this.team.id]);
	}	
}