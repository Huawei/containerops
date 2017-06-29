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

import {
  Directive, HostBinding, Inject, Input, OnInit, OnDestroy
} from '@angular/core';

import { AccordionDirective } from './accordion.directive';

@Directive({
  selector: '[appAccordionLink]'
})
export class AccordionLinkDirective implements OnInit, OnDestroy {

  @Input() public group: any;

  @HostBinding('class.open')
  @Input()
  get open(): boolean {
    return this._open;
  }

  set open(value: boolean) {
    this._open = value;
    if (value) {
      this.nav.closeOtherLinks(this);
    }
  }

  protected _open: boolean;
  protected nav: AccordionDirective;

  public constructor(@Inject(AccordionDirective) nav: AccordionDirective) {
    this.nav = nav;
  }

  public ngOnInit(): any {
    this.nav.addLink(this);
  }

  public ngOnDestroy(): any {
    this.nav.removeGroup(this);
  }

  public toggle(): any {
    this.open = !this.open;
  }
}
