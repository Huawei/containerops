import { ModuleWithProviders }  from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { IndexComponent }  from './index.component';

// org
import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { OrgDetailComponent }  from '../organization/detail.component';
import { OrgAddTeamComponent }  from '../organization/addteam.component';

// team
import { TeamListComponent } from '../team/list.component';
import { TeamAddComponent } from '../team/add.component';
import { TeamDetailComponent } from '../team/detail.component';
import { TeamAddMemberComponent } from '../team/add-member.component';
// repo
import { RepoListComponent }  from '../repository/list.component';

const indexRouting: Routes = [
  {
    path: '',
    component: IndexComponent,
    children: [
      { 
      	path: 'organization', 
      	component: OrgListComponent 
      },
      { 
      	path: 'organization/add', 
      	component: OrgAddComponent 
      },
      { 
        path: 'organization/:id', 
        component: OrgDetailComponent 
      },
      { 
        path: 'organization/:id/addteam', 
        component: OrgAddTeamComponent 
      },
      { 
      	path: 'team', 
      	component: TeamListComponent 
      },
      { 
        path: 'team/add', 
        component: TeamAddComponent 
      },
      { 
        path: 'team/:id', 
        component: TeamDetailComponent
      },
      { 
        path: 'team/:id/member/add', 
        component: TeamAddMemberComponent 
      },
      { 
      	path: 'repository', 
      	component: RepoListComponent 
      }
    ]
  }
];

export const routing: ModuleWithProviders = RouterModule.forChild(indexRouting);