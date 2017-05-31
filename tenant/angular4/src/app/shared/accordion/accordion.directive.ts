import { Directive, OnInit } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';

import { AccordionLinkDirective } from './accordionlink.directive';

import { Subscription } from 'rxjs/Subscription';

@Directive({
  selector: '[appAccordion]',
})
export class AccordionDirective implements OnInit {

  protected navlinks: Array<AccordionLinkDirective> = [];

  private _router: Subscription;

  public closeOtherLinks(openLink: AccordionLinkDirective): void {
    this.navlinks.forEach((link: AccordionLinkDirective) => {
      if (link !== openLink) {
        link.open = false;
      }
    });
  }

  public addLink(link: AccordionLinkDirective): void {
    this.navlinks.push(link);
  }

  public removeGroup(link: AccordionLinkDirective): void {
    const index = this.navlinks.indexOf(link);
    if (index !== -1) {
      this.navlinks.splice(index, 1);
    }
  }

  public getUrl() {
    return this.router.url;
  }

  public ngOnInit(): any {
    this._router = this.router.events.filter(event => event instanceof NavigationEnd).subscribe((event: NavigationEnd) => {
      this.navlinks.forEach((link: AccordionLinkDirective) => {
        if (link.group) {
          const routeUrl = this.getUrl();
          const currentUrl = routeUrl.split('/');
          if (currentUrl.indexOf( link.group ) > 0) {
            link.open = true;
            this.closeOtherLinks(link);
          }
        }
      });
    });
  }

  constructor( private router: Router) {}
}
