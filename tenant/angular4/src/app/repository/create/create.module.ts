/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { CreateComponent } from './create.component';
import { RepoRoutes } from './create.routing';
// import { BrowserModule } from '@angular/platform-browser';
import { MdIconModule, MdInputModule,MdCheckboxModule,MdRadioModule, MdCardModule, MdButtonModule, MdListModule, MdProgressBarModule, MdMenuModule } from '@angular/material';
import { FormsModule }   from '@angular/forms'; 
import { ReactiveFormsModule } from '@angular/forms'
import { MdToolbarModule } from '@angular/material';

import { FileUploadModule } from 'ng2-file-upload/ng2-file-upload';
import { TreeModule } from 'angular-tree-component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';




//import { Repo } from './repo';

@NgModule({
  imports: [
  CommonModule,
  RouterModule.forChild(RepoRoutes),
  MdIconModule,
  MdCardModule,
  FlexLayoutModule,
  MdButtonModule,
  MdListModule,
  MdProgressBarModule,
  MdIconModule,
  MdMenuModule,
  MdCheckboxModule,
  MdRadioModule,
  MdInputModule
  // BrowserModule
   // Repo
    ,FormsModule

   ,ReactiveFormsModule

  ],
  declarations: [CreateComponent],



})

export class RepoModule {}

