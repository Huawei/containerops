import { Component } from '@angular/core';
// import { NavComponent } from './nav.component';
// import { RepoService } from './repo.service';


@Component({
  selector: 'my-app',
  template: `<router-outlet></router-outlet>`,
  // template: require('./app.html'),
  styles: [require('../sass/application.scss')]
})
export class AppComponent { 
	constructor() {}
}