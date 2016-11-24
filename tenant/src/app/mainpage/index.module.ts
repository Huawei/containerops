import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';

import { IndexComponent }  from './index.component';

import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { TeamListComponent } from '../team/list.component';
import { RepoListComponent }  from '../repository/list.component';

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
    TeamListComponent,
  	RepoListComponent
  ],
  providers: [
    // OrgService,
    // RepoService
  ],
  bootstrap: [ IndexComponent ]
})
export class IndexModule { }