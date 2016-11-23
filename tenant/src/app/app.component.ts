import { Component, ViewEncapsulation } from '@angular/core';
// import { NavComponent } from './nav.component';
// import { RepoService } from './repo.service';


@Component({
  selector: 'body',
  template: `<router-outlet></router-outlet>`,
  // template: require('./app.html'),
  styles: [require('../sass/application.scss')],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent { 
	constructor() {}
}