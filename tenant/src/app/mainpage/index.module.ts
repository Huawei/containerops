import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';

import { IndexComponent }  from './index.component';

// org
import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { OrgDetailComponent }  from '../organization/detail.component';
import { OrgAddTeamComponent }  from '../organization/addteam.component';

// team
import { TeamListComponent } from '../team/list.component';
import { TeamAddComponent }  from '../team/add.component';
import { TeamDetailComponent }  from '../team/detail.component';
import { TeamAddMemberComponent } from '../team/add-member.component';

// repo
import { RepoListComponent }  from '../repository/list.component';

// services
import { OrgDataService } from '../organization/data.service';
import { TeamDataService } from '../team/data.service';

import { routing } from './index.router';


@NgModule({
  imports: [ 
  	BrowserModule,
    FormsModule,
    HttpModule,
  	routing
  ],
  declarations: [ 
    IndexComponent,
    OrgListComponent,
    OrgAddComponent,
    OrgDetailComponent,
    OrgAddTeamComponent,
    TeamListComponent,
    TeamAddComponent,
    TeamDetailComponent,
    TeamAddMemberComponent,
  	RepoListComponent
  ],
  providers: [
    OrgDataService,
    TeamDataService
  ],
  bootstrap: [ IndexComponent ]
})
export class IndexModule { }