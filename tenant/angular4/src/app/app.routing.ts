import { Routes } from '@angular/router';

import { CoreLayoutComponent } from './layouts/core/core-layout.component';

export const AppRoutes: Routes = [{
  path: '',
  redirectTo: 'home',
  pathMatch: 'full',
}, {
  path: '',
  component: CoreLayoutComponent,
  children: [{
    path: 'home',
    loadChildren: './dashboard/dashboard.module#DashboardModule'
  }]
}, {
  path: '**',
  redirectTo: 'session/404'
}];
