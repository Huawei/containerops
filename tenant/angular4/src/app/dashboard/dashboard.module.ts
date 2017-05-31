import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MdIconModule, MdCardModule, MdButtonModule, MdListModule, MdProgressBarModule, MdMenuModule } from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';

import { DashboardComponent } from './dashboard.component';
import { DashboardRoutes } from './dashboard.routing';

@NgModule({
  imports: [
    RouterModule.forChild(DashboardRoutes),
    MdIconModule,
    MdCardModule,
    FlexLayoutModule
  ],
  declarations: [ DashboardComponent ]
})

export class DashboardModule {}
