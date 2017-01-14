import { NgModule }      from '@angular/core';
import { BrowserModule, Title } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';

import { IndexModule }  from './mainpage/index.module';

import { AppComponent }  from './app.component';

// import { IndexComponent }  from './mainpage/index.component';
import { LoginComponent }  from './user/login.component';
import { RegisterComponent }  from './user/register.component';

// services
import { UserService }      from './user/user.service';
import { NotifyService }      from './common/notify.service';

import { routing } from './app.router';


@NgModule({
  imports: [ 
  	BrowserModule,
    FormsModule,
    HttpModule,
    IndexModule,
  	routing
  ],
  declarations: [ 
    AppComponent,
    // IndexComponent,
    LoginComponent,
  	RegisterComponent
  ],
  providers: [
    UserService,
    NotifyService,
    Title
  ],
  bootstrap: [ AppComponent ]
})
export class AppModule { }