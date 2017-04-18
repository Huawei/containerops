import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { WorkflowComponent } from './workflow/app.workflow';
import { WorkflowDesignComponent } from './workflow/design/app.workflow.design';

@NgModule({
  declarations: [
    AppComponent,
    WorkflowComponent,
    WorkflowDesignComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
