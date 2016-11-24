import { ModuleWithProviders }  from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { IndexComponent }  from './index.component';

import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { TeamListComponent } from '../team/list.component';
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
      	path: 'team', 
      	component: TeamListComponent 
      },
      { 
      	path: 'repository', 
      	component: RepoListComponent 
      }
    ]
  }
];

export const routing: ModuleWithProviders = RouterModule.forChild(indexRouting);